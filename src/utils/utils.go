package utils

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"golang.org/x/exp/constraints"
)

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m Map[K, V]) Values() []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func Sort[T constraints.Ordered](list []T) []T {
	clone := slices.Clone(list)

	slices.Sort(clone)

	return clone
}

func Read(msg string) (string, error) {
	fmt.Print(msg)

	reader := bufio.NewReader(os.Stdin)

	res, err := reader.ReadString('\n')

	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(res), nil
}
