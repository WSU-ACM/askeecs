package kvstore

var stores = map[string]store{}

type store map[string]interface{}

func Set(store_name, key string, i interface{}) bool {
	s, found := stores[store_name]

	if !found {
		s = make(store)
		stores[store_name] = s
	}

	s[key] = i

	return true
}

func Get(store_name, key string) (interface{}, bool) {
	var s     store
	var v     interface{}
	var found bool

	s, found = stores[store_name]

	if !found {
		return nil, false
	}

	v, found = s[key]

	return v, found
}
