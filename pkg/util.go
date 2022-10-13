package pkg

import "net/netip"

func getAddrRange(prefix netip.Prefix) (netip.Addr, netip.Addr){
	start := prefix.Addr()
	end := start

	for {
		if (prefix.Contains(end.Next()) == false){
			return start,end
		}
		end = end.Next()
	}
} 