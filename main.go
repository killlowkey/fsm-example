package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/looplab/fsm"
)

// 定义状态常量
const (
	StateCreated   = "created"
	StatePaid      = "paid"
	StateShipped   = "shipped"
	StateCompleted = "completed"
)

// 定义事件常量
const (
	EventPay      = "pay"
	EventShip     = "ship"
	EventComplete = "complete"
)

// Order 定义一个订单结构体
type Order struct {
	fsm *fsm.FSM
}

// NewOrder 创建一个新的订单，并初始化状态机
func NewOrder() *Order {
	order := &Order{
		fsm: fsm.NewFSM(
			StateCreated, // 初始状态
			fsm.Events{
				// 状态转换定义
				{Name: EventPay, Src: []string{StateCreated}, Dst: StatePaid},
				{Name: EventShip, Src: []string{StatePaid}, Dst: StateShipped},
				{Name: EventComplete, Src: []string{StateShipped}, Dst: StateCompleted},
			},
			fsm.Callbacks{
				// 在任何事件触发之前调用
				"before_event": func(ctx context.Context, e *fsm.Event) {
					fmt.Printf("[before_event] Event '%s' about to trigger from state '%s'\n", e.Event, e.Src)
				},
				// 在任何事件触发之后调用
				"after_event": func(ctx context.Context, e *fsm.Event) {
					fmt.Printf("[after_event] Event '%s' triggered, current state is '%s'\n", e.Event, e.Dst)
				},
				// 在离开任何状态时调用
				"leave_state": func(ctx context.Context, e *fsm.Event) {
					fmt.Printf("[leave_state] Leaving state '%s'\n", e.Src)
				},
				// 在进入任何状态时调用
				"enter_state": func(ctx context.Context, e *fsm.Event) {
					fmt.Printf("[enter_state] Entered state '%s'\n", e.Dst)
				},
				// 在事件 `pay` 之前调用，执行业务相关的逻辑
				"before_pay": func(ctx context.Context, e *fsm.Event) {
					fmt.Println("[before_pay] Checking if payment can proceed...")
					// 支付相关逻辑，检查是否已经支付
					if len(e.Args) < 2 {
						fmt.Println("[before_pay] Insufficient arguments, payment canceled.")
						e.Cancel(errors.New("invalid payment details"))
					} else {
						fmt.Println("[before_pay] Payment details are valid.")
					}
				},
				// 在事件 `pay` 之后调用
				"after_pay": func(ctx context.Context, e *fsm.Event) {
					fmt.Println("[after_pay] Payment done!")
				},
				// 在事件 `ship` 之前调用
				"before_ship": func(ctx context.Context, e *fsm.Event) {
					fmt.Println("[before_ship] Checking if shipment is ready...")
					if e.Src != StatePaid {
						fmt.Println("[before_ship] Cannot ship, order must be paid first!")
						e.Cancel(errors.New("order not paid yet"))
					}
				},
				// 在事件 `complete` 之前调用
				"before_complete": func(ctx context.Context, e *fsm.Event) {
					fmt.Println("[before_complete] Checking if order can be completed...")
					if e.Src != StateShipped {
						fmt.Println("[before_complete] Cannot complete order, it must be shipped first!")
						e.Cancel(errors.New("order not shipped yet"))
					}
				},
			},
		),
	}
	return order
}

// CurrentState 返回订单的当前状态
func (o *Order) CurrentState() string {
	return o.fsm.Current()
}

// Pay 执行支付操作
// before_pay -> before_event -> leave_state -> enter_state -> after_pay -> after_event
func (o *Order) Pay() error {
	return o.fsm.Event(context.TODO(), EventPay, "order1", "order2")
}

// Ship 执行发货操作
func (o *Order) Ship() error {
	return o.fsm.Event(context.TODO(), EventShip)
}

// Complete 完成订单
func (o *Order) Complete() error {
	return o.fsm.Event(context.TODO(), EventComplete)
}

func main() {
	order := NewOrder()
	fmt.Println("Initial State:", order.CurrentState())

	// 尝试支付订单
	if err := order.Pay(); err != nil {
		fmt.Println("Error paying:", err)
	} else {
		fmt.Println("Order paid. Current State:", order.CurrentState())
	}

	// 尝试发货订单
	if err := order.Ship(); err != nil {
		fmt.Println("Error shipping:", err)
	} else {
		fmt.Println("Order shipped. Current State:", order.CurrentState())
	}

	// 尝试完成订单
	if err := order.Complete(); err != nil {
		fmt.Println("Error completing:", err)
	} else {
		fmt.Println("Order completed. Current State:", order.CurrentState())
	}
}
