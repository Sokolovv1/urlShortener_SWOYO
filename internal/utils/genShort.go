package utils

func GenShort(id int) string {
	var result string
	for id > 0 {
		id-- // Декрементируем для правильного отображения букв
		result = string(rune('A'+id%26)) + result
		id /= 26
	}
	return result
}
