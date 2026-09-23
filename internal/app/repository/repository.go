package repository

import "fmt"

type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// AlgoBenchmark - бенчмарк производительности алгоритма
type AlgoBenchmark struct {
	BenchmarkID     string
	AlgoName        string
	Complexity      string
	InputSize       int    // Объём данных (числовое поле для фильтра)
	ExecTimeMs      int    // Время выполнения (числовое поле)
	MemoryMb        float64 // Потребление памяти (доп. числовое поле для темы)
	Category        string
	PreviewImageKey string
	VideoKey        string
	Status          string // published, draft, deleted
	Likes           []string
	Description     string
}

func (r *Repository) GetBenchmarks() ([]AlgoBenchmark, error) {
	benchmarks := []AlgoBenchmark{
		{
			BenchmarkID: "bench_001", AlgoName: "Сортировка пузырьком", Complexity: "O(n²)",
			InputSize: 1000, ExecTimeMs: 150, MemoryMb: 2.1, Category: "сортировка",
			PreviewImageKey: "bubble.png", VideoKey: "algorithm_sorting_on_squared.mp4",
			Status: "published", Likes: []string{"u1", "u2", "u3"},
			Description: "Простой, но медленный алгоритм. Квадратичная сложность делает его неэффективным для больших массивов.",
		},
		{
			BenchmarkID: "bench_002", AlgoName: "Бинарный поиск", Complexity: "O(log n)",
			InputSize: 1000000, ExecTimeMs: 5, MemoryMb: 0.5, Category: "поиск",
			PreviewImageKey: "binary.png", VideoKey: "search_binary_search_ologn.mp4",
			Status: "published", Likes: []string{"u1"},
			Description: "Молниеносный поиск в отсортированном массиве. Делит область поиска пополам на каждом шаге.",
		},
		{
			BenchmarkID: "bench_003", AlgoName: "Прямой доступ", Complexity: "O(1)",
			InputSize: 1, ExecTimeMs: 1, MemoryMb: 0.1, Category: "доступ",
			PreviewImageKey: "constant.png", VideoKey: "big_o_notation_overview.mp4",
			Status: "published", Likes: []string{"u2", "u3"},
			Description: "Мгновенное получение элемента по индексу. Время выполнения не зависит от размера данных.",
		},
		{
			BenchmarkID: "bench_004", AlgoName: "Быстрая сортировка", Complexity: "O(n log n)",
			InputSize: 5000, ExecTimeMs: 12, MemoryMb: 3.5, Category: "сортировка",
			PreviewImageKey: "quick_sort.png", VideoKey: "quick_sort.mp4",
			Status: "published", Likes: []string{"u1", "u2", "u3", "u4"},
			Description: "Один из самых быстрых алгоритмов на практике. Использует стратегию 'разделяй и властвуй'.",
		},
		{
			BenchmarkID: "bench_005", AlgoName: "Линейный поиск", Complexity: "O(n)",
			InputSize: 500, ExecTimeMs: 3, MemoryMb: 0.2, Category: "поиск",
			PreviewImageKey: "linear.png", VideoKey: "linear.mp4",
			Status: "published", Likes: []string{"u2"},
			Description: "Последовательная проверка каждого элемента. Подходит только для небольших или несортированных данных.",
		},
		{
			BenchmarkID: "bench_006", AlgoName: "Сортировка слиянием", Complexity: "O(n log n)",
			InputSize: 8000, ExecTimeMs: 18, MemoryMb: 8.0, Category: "сортировка",
			PreviewImageKey: "merge.png", VideoKey: "merge.mp4",
			Status: "published", Likes: []string{"u1", "u3", "u5"},
			Description: "Стабильная сортировка. Требует дополнительной памяти, но гарантирует O(n log n) в худшем случае.",
		},
		{
			BenchmarkID: "bench_007", AlgoName: "Хеш-таблица", Complexity: "O(1)",
			InputSize: 1, ExecTimeMs: 1, MemoryMb: 4.0, Category: "хеш",
			PreviewImageKey: "hash.png", VideoKey: "hash.mp4",
			Status: "draft", Likes: []string{},
			Description: "Черновик бенчмарка. Мгновенный поиск по ключу, но с высоким потреблением памяти.",
		},
		{
			BenchmarkID: "bench_008", AlgoName: "Фибоначчи (рекурсия)", Complexity: "O(2^n)",
			InputSize: 20, ExecTimeMs: 5000, MemoryMb: 15.0, Category: "рекурсия",
			PreviewImageKey: "bubble.png", VideoKey: "bubble.mp4",
			Status: "deleted", Likes: []string{"u1"},
			Description: "Удалённый бенчмарк. Экспоненциальный рост времени делает этот подход непригодным для n > 30.",
		},
	}

	if len(benchmarks) == 0 {
		return nil, fmt.Errorf("коллекция пуста")
	}
	return benchmarks, nil
}

func (r *Repository) GetBenchmarkByID(id string) (AlgoBenchmark, error) {
	benchmarks, err := r.GetBenchmarks()
	if err != nil {
		return AlgoBenchmark{}, err
	}
	for _, b := range benchmarks {
		if b.BenchmarkID == id {
			return b, nil
		}
	}
	return AlgoBenchmark{}, fmt.Errorf("бенчмарк не найден")
}

func (r *Repository) GetDraft() (AlgoBenchmark, error) {
	benchmarks, err := r.GetBenchmarks()
	if err != nil {
		return AlgoBenchmark{}, err
	}
	for _, b := range benchmarks {
		if b.Status == "draft" {
			return b, nil
		}
	}
	return AlgoBenchmark{}, fmt.Errorf("черновик не найден")
}

func GetMinioURL(key string) string {
	return fmt.Sprintf("http://localhost:9000/algorithms-media/%s", key)
}