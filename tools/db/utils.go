package main

import (
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"sort"
)

func rnd(min, max uint64) uint64 {
	return rand.Uint64()%(max-min+1) + min
}

func flip(success float32) bool {
	return rand.Float32() < success
}

func randomIPv4() net.IP {
	ip := make([]byte, 4)
	rand.Read(ip)
	return net.IP(ip)
}

// IPv4Subnet gets the ipv4 subnet given the ip
func IPv4Subnet(ip net.IP) *net.IPNet {
	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(24, 32),
	}
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func insertQuery(v interface{}) string {
	query := fmt.Sprintf("INSERT INTO %s (", reflect.Indirect(reflect.ValueOf(v)).Type().Name())
	vals := "("
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		if i == 0 {
			v := fmt.Sprint(val.Field(i))
			if v == "" {
				v = "NULL"
			}
			if v != "NULL" && val.Field(i).Type().Name() == "string" {
				v = fmt.Sprintf(`'%s'`, v)
			}
			query = fmt.Sprintf("%s%s", query, val.Type().Field(i).Name)
			vals = fmt.Sprintf("%s%s", vals, v)
		} else {
			v := fmt.Sprint(val.Field(i))
			if v == "" {
				v = "NULL"
			}
			if v != "NULL" && val.Field(i).Type().Name() == "string" {
				v = fmt.Sprintf(`'%s'`, v)
			}
			query = fmt.Sprintf("%s, %s", query, val.Type().Field(i).Name)
			vals = fmt.Sprintf("%s, %s", vals, v)
		}
	}
	query = fmt.Sprintf("%s) VALUES %s);", query, vals)
	return query
}

func randomKey(m map[uint64]struct{}) uint64 {
	keys := []uint64{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys[rnd(0, uint64(len(keys)-1))]
}
