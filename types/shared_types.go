package types

type Job struct {
	Prompt         string
	Model          string
	Id             string
	MaxAttempts    int
	CurentAttempts int
}

type MealParts struct {
	MajorParts []string `json:"meal_parts"`
}

type SymanticSearchResult struct {
	FoodName string        `json:"foodName"`
	FoodId   string        `json:"foodId"`
	Servings []ServingInfo `json:"servings"`
}

type ServingInfo struct {
	ServingId     string `json:"servingId"`
	ServingSize   string `json:"servingSize"`
	ServingWeight string `json:"servingWeight"`
}

type ServingSizeChoices struct {
	FoodName           string `json:"foodName"`
	RequestedServingId string `json:"requestedServingId"`
}

type NutritionalInfo struct {
	FoodName      string     `json:"foodName"`
	ServingSize   string     `json:"servingSize"`
	ServingWeight string     `json:"servingWeight"`
	Nutrients     []Nutrient `json:"nutrients"`
}

type Nutrient struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
