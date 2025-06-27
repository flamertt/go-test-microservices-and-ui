package model

type Book struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Publisher    string `json:"publisher"`
	Author       string `json:"author"`
	CategoryName string `json:"category_name"`
	ProductCode  string `json:"product_code"`
	PageCount    int    `json:"page_count"`
	ReleasedYear int    `json:"released_year"`
}

type Author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Biography string `json:"biography"`
}

type Genre struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Recommendation struct {
	Book   Book   `json:"book"`
	Reason string `json:"reason"`
	Score  int    `json:"score"`
}

type RecommendationResponse struct {
	Recommendations []Recommendation `json:"recommendations"`
	Total           int              `json:"total"`
	Category        string           `json:"category,omitempty"`
	Author          string           `json:"author,omitempty"`
	Timestamp       string           `json:"timestamp"`
}

// API Response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
} 