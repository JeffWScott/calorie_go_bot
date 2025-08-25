package pipeline

import (
	"calorie_bot/api"
	"calorie_bot/db"
	"calorie_bot/types"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Pipeline struct {
	db  *db.Db
	api *api.OpenRouter
	job types.Job
}

func New(api *api.OpenRouter, db *db.Db) *Pipeline {

	p := &Pipeline{
		db:  db,
		api: api,
	}

	return p
}

func (p *Pipeline) Log(message string) {
	fmt.Printf("[%s] %s\n", p.job.Id, message)
}

func (p *Pipeline) Run(job types.Job) (string, error) {
	p.job = job
	p.Log("Starting New Job")
	p.Log(fmt.Sprintf("Prompt: %s", p.job.Prompt))

	result, err := p.Start()
	if err != nil {
		p.Log(fmt.Sprintf("job failed: %v", err))
		return "", fmt.Errorf("job failed: %v", err)
	}

	return result, err

}

func (p *Pipeline) Start() (string, error) {
	/**
	GET MEAL PARTS
	*/
	p.job.MaxAttempts = 3
	p.job.CurentAttempts = 0
	mealParts, err := p.GetMealParts()
	if err != nil {
		return "", err
	}
	p.Log(fmt.Sprintf("mealParts: %v", mealParts))

	/**
	SCHEDULE GET NUTRITION
	  - offshoot process
	*/
	go p.ScheduleGetNutrition(mealParts)

	/**
	GET SYMANTIC SEARCH RESULTS
	*/
	p.job.MaxAttempts = 3
	p.job.CurentAttempts = 0
	symanticSearch, err := p.GetSymanticSearchResults(mealParts)
	if err != nil {
		return "", err
	}

	prettyJSON, err := json.MarshalIndent(symanticSearch, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	fmt.Println(string(prettyJSON))

	/**
	PICK SERVINGS
	*/

	p.job.MaxAttempts = 3
	p.job.CurentAttempts = 0
	servingSizeChoices, err := p.GetServingSizes(symanticSearch)
	if err != nil {
		return "", err
	}

	prettyJSON, err = json.MarshalIndent(servingSizeChoices, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(prettyJSON))

	/**
	GET NUTRITIONAL INFO
	*/
	p.job.MaxAttempts = 3
	p.job.CurentAttempts = 0
	nutritionalInfo, err := p.GetNutritionalInfo(servingSizeChoices)
	if err != nil {
		return "", err
	}
	/*
		prettyJSON, err = json.MarshalIndent(nutritionalInfo, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}
		fmt.Println(string(prettyJSON))
	*/

	/**
	PICK SERVINGS
	*/

	p.job.MaxAttempts = 3
	p.job.CurentAttempts = 0
	mealMacros, err := p.GetMealMacros(nutritionalInfo)
	if err != nil {
		return "", err
	}

	fmt.Println(mealMacros)

	return mealMacros, nil
}

func (p *Pipeline) GetMealParts() (types.MealParts, error) {
	p.job.CurentAttempts++
	if p.job.CurentAttempts > p.job.MaxAttempts {
		return types.MealParts{}, errors.New("[GetMealParts] reached max attempts")
	}
	p.Log(fmt.Sprintf("[GetMealParts] Running %d/%d", p.job.CurentAttempts, p.job.MaxAttempts))

	client := p.api.NewClient()

	mealParts, err := client.GetMealParts(p.job.Prompt, p.job.Id)
	if err != nil {
		p.Log(fmt.Sprintf("[GetMealParts] Error: %v", err))
		p.GetMealParts()
	}

	return mealParts, nil
}

func (p *Pipeline) ScheduleGetNutrition(mealParts types.MealParts) {
	p.Log(fmt.Sprintf("[GetSymanticSearchResults] Running %d/%d", p.job.CurentAttempts, p.job.MaxAttempts))

	convexClinet := p.db.NewClient()

	err := convexClinet.ScheduleGetNutrition(mealParts.MajorParts, p.job.Id)
	if err != nil {
		p.Log(fmt.Sprintf("[ScheduleGetNutrition] Error: %v", err))
	}
}

func (p *Pipeline) GetSymanticSearchResults(mealParts types.MealParts) ([]types.SymanticSearchResult, error) {
	p.job.CurentAttempts++
	if p.job.CurentAttempts > p.job.MaxAttempts {
		return []types.SymanticSearchResult{}, errors.New("[GetSymanticSearchResults] reached max attempts")
	}
	p.Log(fmt.Sprintf("[GetSymanticSearchResults] Running %d/%d", p.job.CurentAttempts, p.job.MaxAttempts))

	convexClinet := p.db.NewClient()

	symanticSearchResults, err := convexClinet.GetSymanticSearchResults(mealParts.MajorParts, p.job.Id)
	if err != nil {
		p.Log(fmt.Sprintf("[GetSymanticSearchResults] Error: %v", err))
		p.GetSymanticSearchResults(mealParts)
	}

	return symanticSearchResults, nil
}

func (p *Pipeline) GetServingSizes(semanticSearchResults []types.SymanticSearchResult) ([]types.ServingSizeChoices, error) {
	p.job.CurentAttempts++
	if p.job.CurentAttempts > p.job.MaxAttempts {
		return []types.ServingSizeChoices{}, errors.New("[GetServingSizes] reached max attempts")
	}
	p.Log(fmt.Sprintf("[GetServingSizes] Running %d/%d", p.job.CurentAttempts, p.job.MaxAttempts))

	client := p.api.NewClient()

	servingSizeChoices, err := client.GetServingSizes(semanticSearchResults, p.job.Prompt, p.job.Id)
	if err != nil {
		p.Log(fmt.Sprintf("[GetServingSizes] Error: %v", err))
		p.GetServingSizes(semanticSearchResults)
	}

	return servingSizeChoices, nil
}

func (p *Pipeline) GetNutritionalInfo(servingSizeChoices []types.ServingSizeChoices) ([]types.NutritionalInfo, error) {
	p.job.CurentAttempts++
	if p.job.CurentAttempts > p.job.MaxAttempts {
		return []types.NutritionalInfo{}, errors.New("[GetNutritionalInfo] reached max attempts")
	}
	p.Log(fmt.Sprintf("[GetNutritionalInfo] Running %d/%d", p.job.CurentAttempts, p.job.MaxAttempts))

	convexClinet := p.db.NewClient()

	nutritionalInfoResults, err := convexClinet.GetNutritionalInfo(servingSizeChoices, p.job.Id)
	if err != nil {
		p.Log(fmt.Sprintf("[GetNutritionalInfo] Error: %v", err))
		p.GetNutritionalInfo(servingSizeChoices)
	}

	return nutritionalInfoResults, nil
}

func (p *Pipeline) GetMealMacros(nutritionalInfo []types.NutritionalInfo) (string, error) {
	p.job.CurentAttempts++
	if p.job.CurentAttempts > p.job.MaxAttempts {
		return "", errors.New("[GetServingSizes] reached max attempts")
	}
	p.Log(fmt.Sprintf("[GetServingSizes] Running %d/%d", p.job.CurentAttempts, p.job.MaxAttempts))

	client := p.api.NewClient()

	mealMacros, err := client.GetMealMacros(p.job.Prompt, nutritionalInfo, p.job.Model, p.job.Id)
	if err != nil {
		p.Log(fmt.Sprintf("[GetServingSizes] Error: %v", err))
		p.GetMealMacros(nutritionalInfo)
	}

	return mealMacros, nil
}
