package cache

type GPTMessageCache interface {
	Set(key string, value string)
	Get(key string) string
}
