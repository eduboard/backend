package url

import "net/url"

func StringifyURLs(urls []url.URL) []string {
	s := make([]string, len(urls))
	for k, v := range urls {
		s[k] = v.String()
	}
	return s
}

func URLifyStrings(strings []string) (urls []url.URL, err error) {
	URLs := make([]url.URL, len(strings))
	for k, v := range strings {
		u, err := url.Parse(v)
		if err != nil {
			return []url.URL{}, err
		}
		URLs[k] = *u
	}
	return URLs, nil
}
