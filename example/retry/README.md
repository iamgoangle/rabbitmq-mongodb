### Step Setup

1. Create **WorkQueue**
2. Create **RetryQueue**
   - define **x-dead-letter-exchange** to **WorkExchange**
   - define **x-message-ttl** to `up-to-you-value`
3. Create **WorkExchange** for queue routing
4. Create **RetryExchange** for retry queue routing

### Dead Letter Exchange

https://www.rabbitmq.com/dlx.html

### Reference

- https://medium.com/@igkuz/ruby-retry-scheduled-tasks-with-dead-letter-exchange-in-rabbitmq-9e38aa39089b
- https://medium.com/@kiennguyen88/rabbitmq-delay-retry-schedule-with-dead-letter-exchange-31fb25a440fc
- https://m.alphasights.com/exponential-backoff-with-rabbitmq-78386b9bec81
