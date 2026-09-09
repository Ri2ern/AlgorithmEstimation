package repository

import "fmt"

type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type AlgorithmEstimation struct {
	EstimationID          string
	ComplexityClass       string
	InputDataVolume       int
	ExecutionTimeMs       int
	OperationCategory     string
	PreviewImageKey       string
	DemonstrationVideoKey string
	PublicationStatus     string
	ApproverIDs           []string
	EstimationTitle       string
	EstimationDescription string
}

func (r *Repository) GetEstimations() ([]AlgorithmEstimation, error) {
	estimations := []AlgorithmEstimation{
		{
			EstimationID: "est_001", ComplexityClass: "O(n²)", InputDataVolume: 1000, ExecutionTimeMs: 150, OperationCategory: "сортировка",
			PreviewImageKey: "bubble.png", DemonstrationVideoKey: "algorithm_sorting_on_squared.mp4",
			PublicationStatus: "published", ApproverIDs: []string{"user1", "user2", "user3"},
			EstimationTitle: "Сортировка пузырьком", EstimationDescription: "Квадратичная сложность. Медленно работает на больших объёмах данных.",
		},
		{
			EstimationID: "est_002", ComplexityClass: "O(log n)", InputDataVolume: 1000000, ExecutionTimeMs: 5, OperationCategory: "поиск",
			PreviewImageKey: "binary.png", DemonstrationVideoKey: "search_binary_search_ologn.mp4",
			PublicationStatus: "published", ApproverIDs: []string{"user1"},
			EstimationTitle: "Бинарный поиск", EstimationDescription: "Логарифмическая сложность. Очень быстрый поиск в отсортированном массиве.",
		},
		{
			EstimationID: "est_003", ComplexityClass: "O(1)", InputDataVolume: 1, ExecutionTimeMs: 1, OperationCategory: "доступ",
			PreviewImageKey: "constant.png", DemonstrationVideoKey: "big_o_notation_overview.mp4",
			PublicationStatus: "published", ApproverIDs: []string{"user2", "user3"},
			EstimationTitle: "Константное время", EstimationDescription: "Мгновенное выполнение. Доступ к элементу массива по индексу.",
		},
		{
			EstimationID: "est_004", ComplexityClass: "O(n log n)", InputDataVolume: 5000, ExecutionTimeMs: 12, OperationCategory: "сортировка",
			PreviewImageKey: "quick_sort.png", DemonstrationVideoKey: "quick_sort.mp4",
			PublicationStatus: "published", ApproverIDs: []string{"user1", "user2", "user3", "user4"},
			EstimationTitle: "Быстрая сортировка", EstimationDescription: "Средняя сложность O(n log n). Быстрый алгоритм разделяй и властвуй.",
		},
		{
			EstimationID: "est_005", ComplexityClass: "O(n)", InputDataVolume: 500, ExecutionTimeMs: 3, OperationCategory: "поиск",
			PreviewImageKey: "linear.png", DemonstrationVideoKey: "linear.mp4",
			PublicationStatus: "published", ApproverIDs: []string{"user2"},
			EstimationTitle: "Линейный поиск", EstimationDescription: "Линейная сложность. Последовательный поиск по массиву.",
		},
		{
			EstimationID: "est_006", ComplexityClass: "O(n log n)", InputDataVolume: 8000, ExecutionTimeMs: 18, OperationCategory: "сортировка",
			PreviewImageKey: "merge.png", DemonstrationVideoKey: "merge.mp4",
			PublicationStatus: "published", ApproverIDs: []string{"user1", "user3", "user5"},
			EstimationTitle: "Сортировка слиянием", EstimationDescription: "Стабильная сортировка O(n log n). Подход разделяй и властвуй.",
		},
		{
			EstimationID: "est_007", ComplexityClass: "O(1)", InputDataVolume: 1, ExecutionTimeMs: 1, OperationCategory: "хеш",
			PreviewImageKey: "hash.png", DemonstrationVideoKey: "hash.mp4",
			PublicationStatus: "draft", ApproverIDs: []string{},
			EstimationTitle: "Хеш-поиск", EstimationDescription: "Черновая оценка. Поиск в хеш-таблице за константное время.",
		},
		{
			EstimationID: "est_008", ComplexityClass: "O(2^n)", InputDataVolume: 20, ExecutionTimeMs: 5000, OperationCategory: "рекурсия",
			PreviewImageKey: "bubble.png", DemonstrationVideoKey: "bubble.mp4",
			PublicationStatus: "deleted", ApproverIDs: []string{"user1"},
			EstimationTitle: "Числа Фибоначчи", EstimationDescription: "Экспоненциальная сложность. Удалённая оценка.",
		},
	}

	if len(estimations) == 0 {
		return nil, fmt.Errorf("коллекция пуста")
	}

	return estimations, nil
}

func (r *Repository) GetEstimationByID(id string) (AlgorithmEstimation, error) {
	estimations, err := r.GetEstimations()
	if err != nil {
		return AlgorithmEstimation{}, err
	}

	for _, est := range estimations {
		if est.EstimationID == id {
			return est, nil
		}
	}

	return AlgorithmEstimation{}, fmt.Errorf("оценка не найдена")
}

func (r *Repository) GetDraft() (AlgorithmEstimation, error) {
	estimations, err := r.GetEstimations()
	if err != nil {
		return AlgorithmEstimation{}, err
	}

	for _, est := range estimations {
		if est.PublicationStatus == "draft" {
			return est, nil
		}
	}

	return AlgorithmEstimation{}, fmt.Errorf("черновик не найден")
}

func GetMinioURL(key string) string {
	return fmt.Sprintf("http://localhost:9000/algorithms-media/%s", key)
}