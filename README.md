## 常见的回调函数类型及其含义
1. `before_<event>`
   1. 含义：在某个特定事件触发之前调用。 
   2. 使用场景：你可以在这个回调中进行检查或准备工作，决定是否允许事件发生。如果逻辑不满足，可以取消事件。
   3. 例子：before_pay：检查订单是否可以支付（比如是否库存足够或支付信息是否正确）。
2. `after_<event>`
   1. 含义：在某个特定事件触发之后调用。 
   2. 使用场景：通常用于在事件执行后做一些清理操作或记录日志。事件已经执行成功，可以在此回调中进行后续处理。 
   3. 例子：after_pay：记录支付完成日志或触发下一步操作，如通知用户支付成功。
3. `leave_<state>`
   1. 含义：当即将离开某个状态时调用。
   2. 使用场景：在状态转换之前，你可以在 leave_<state> 回调中做一些必要的退出操作，比如保存状态数据或验证是否允许离开该状态。
   3. 例子：leave_paid：在订单从“已支付”状态离开时，可能需要记录发票生成等操作。
4. `enter_<state>`
   1. 含义：在进入某个状态时调用。
   2. 使用场景：通常用于在进入某个状态时执行初始化操作，例如分配资源、记录状态变更、发送通知等。
   3. 例子：enter_shipped：在订单进入“已发货”状态时，可以自动触发发货通知给用户或更新物流信息。
5. `before_event`
   1. 含义：在任何事件触发之前调用。
   2. 使用场景：你可以使用 before_event 来对所有事件做统一的预处理，比如日志记录或全局的权限检查。
   3. 例子：记录所有事件触发的尝试，并且可以统一处理权限控制。
6. `after_event`
   1. 含义：在任何事件触发之后调用。
   2. 使用场景：用于在任何事件完成后做一些全局性的清理或记录。可以进行全局监控或记录事件结果。
   3. 例子：事件执行后，记录每个状态转换的详细日志，便于后续分析。
7. `enter_state`
   1. 含义：在进入任何状态时调用。
   2. 使用场景：这个回调是在进入任意状态时触发，用于做一些全局性的处理。适合场景如全局状态初始化、全局通知等。
   3. 例子：在系统进入任意状态时，都可以在此统一发送通知或者更新系统状态。
8. `leave_state`
   1. 含义：在离开任何状态时调用。
   2. 使用场景：用于在离开任意状态时触发清理或退出逻辑。例如，释放某些资源，或者在状态变更前保存一些系统快照。
   3. 例子：在状态离开时自动保存状态日志，防止数据丢失。

## 回调函数的执行顺序
当一个事件被触发时，状态机会依次执行以下回调函数（如果定义了这些回调）：

1. `before_<event>`：首先检查是否可以触发事件。 
2. `leave_<state>`：如果事件允许，会从当前状态离开，并执行该回调。 
3. `enter_<state>`：接着进入目标状态，并执行该回调。 
4. `after_<event>`：最后执行事件成功后的回调。

## 程序输出
成功处理
```text
Initial State: created
[before_pay] Checking if payment can proceed...
[before_pay] Payment details are valid.
[before_event] Event 'pay' about to trigger from state 'created'
[leave_state] Leaving state 'created'
[enter_state] Entered state 'paid'
[after_pay] Payment done!
[after_event] Event 'pay' triggered, current state is 'paid'
Order paid. Current State: paid
[before_ship] Checking if shipment is ready...
[before_event] Event 'ship' about to trigger from state 'paid'
[leave_state] Leaving state 'paid'
[enter_state] Entered state 'shipped'
[after_event] Event 'ship' triggered, current state is 'shipped'
Order shipped. Current State: shipped
[before_complete] Checking if order can be completed...
[before_event] Event 'complete' about to trigger from state 'shipped'
[leave_state] Leaving state 'shipped'
[enter_state] Entered state 'completed'
[after_event] Event 'complete' triggered, current state is 'completed'
Order completed. Current State: completed
```

失败处理
```text
Initial State: created
[before_pay] Checking if payment can proceed...
[before_pay] Insufficient arguments, payment canceled.
Error paying: transition canceled with error: invalid payment details
Error shipping: event ship inappropriate in current state created
Error completing: event complete inappropriate in current state created
```