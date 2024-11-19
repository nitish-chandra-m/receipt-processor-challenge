package models

type Rubric map[string]float64

// NewRubric initializes a Rubric with the provided values or defaults.
func NewRubric(rubric map[string]float64) Rubric {
	if len(rubric) == 0 {
		// Provide a default rubric
		return Rubric{
			"retailer":          1,
			"roundedTotal":      50,
			"totalMultiple":     25,
			"pairOfItems":       5,
			"descriptionLength": 0.2,
			"oddPurchaseDay":    6,
			"afternoonPurchase": 10,
		}
	}
	return Rubric(rubric)
}
