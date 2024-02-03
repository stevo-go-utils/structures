# structures
Unique datastructures in go.

## How to install
```go get https://github.com/stevo-go-utils/structures```

## Balancer
**What is it for?**
The best use case for this is to rotate proxies with the advantage to space them out easily and have the ability to remove bad proxies from the list.

**How does rotation work?**
The data structure for rotation is a circular queue. This means that its just like a line but instead of popping upon use the proxy that is popped goes back to the end of the line. This assures that we are spreading the use of proxies completely.

**What are the custom functions I added?**
There are three utility operations that I added to responses. Use() indicates to the stats for the balancer the last `time.Time` it was used. This is useful for when the timeout is set in the balancer. Report() indicates to the balancer that this proxy was the reason for an error and it will uptick its count of error reports in the statistics. If the errors option is set upon balancer creation this allows the balancer to delete the proxy from the list if the errors limit was reached. If no option is set, then it would auto-delete. Wait() will run a `time.Sleep(REMAINING_TIMEOUT)`. This only happens if the timeout option is set. This allows us to make sure we are at least spacing out the uses by X timeout time.

**Type Explanations**
The V comparable allows us to rotate any comparable datatype (thanks to go generics).
```go
type Balancer[V comparable] struct {
    cll   CircularLinkedList[V]
    stats *SafeMap[V, *BalancerStats]
    *BalancerOpts
}

type BalancerStats struct {
    errors   int // Num of error reports (hitting past limit will delete item)
    lastUsed time.Time // Last time data was used 
    // (for making sure that uses are guarenteed to be spaced out by X time)
}

type BalancerResp[V comparable] struct {
    Data   func() V // Get resp data
    Use    func() // Update last used time in stats
    Report func() // Uptick error count in stats
    Wait   func() // Wait remaining timeout
}
```

**Options**
```go
type BalancerOpts struct {
    MaxErrs    int
    UseTimeout *time.Duration
}

func DefaultBalancerOpts() *BalancerOpts {
    return &BalancerOpts{
        MaxErrs:    -1,
        UseTimeout: nil,
    }
}

  

func MaxErrsBalancerOpt(maxErrs int) BalancerOpt {
    return func(opts *BalancerOpts) {
        opts.MaxErrs = maxErrs
    }
}

  

func UseTimeoutBalancerOpt(useTimeout time.Duration) BalancerOpt {
    return func(opts *BalancerOpts) {
        opts.UseTimeout = &useTimeout
    }
}
```
