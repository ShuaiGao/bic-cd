package limiter

// Limiter 限流器
// 限流器用于限制本地请求第三方接口流量限流
// 采用最严格限流
type Limiter interface {
	// Run 限流执行方法 f
	Run(f func())
	// RunMulti 限流执行方法 f
	// 并指定当前方法的优先级，占用单次流量时间倍数
	RunMulti(priority int, multiple float32, f func())
	// Quit 退出限流器
	// 正在限流中的方法 f 会被丢弃
	Quit()
}
