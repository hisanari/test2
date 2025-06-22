package mermaid

import (
	"testing"
)

func TestParseERDiagram(t *testing.T) {
	input := `erDiagram
    USER {
        int id PK
        string name
        string email UK
    }
    ORDER {
        int id PK
        int user_id FK
        decimal total
    }
    USER ||--o{ ORDER : places`
	
	diagram, err := ParseERDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Check entities
	if len(diagram.Entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(diagram.Entities))
	}
	
	// Check USER entity
	userEntity := diagram.Entities[0]
	if userEntity.Name != "USER" {
		t.Errorf("Expected entity name 'USER', got '%s'", userEntity.Name)
	}
	
	if len(userEntity.Attributes) != 3 {
		t.Errorf("Expected 3 attributes for USER, got %d", len(userEntity.Attributes))
	}
	
	// Check primary key attribute
	idAttr := userEntity.Attributes[0]
	if !idAttr.IsPK {
		t.Error("Expected id attribute to be primary key")
	}
	
	// Check unique constraint
	emailAttr := userEntity.Attributes[2]
	if !emailAttr.IsUnique {
		t.Error("Expected email attribute to be unique")
	}
	
	// Check relationships
	if len(diagram.Relationships) != 1 {
		t.Errorf("Expected 1 relationship, got %d", len(diagram.Relationships))
	}
	
	rel := diagram.Relationships[0]
	if rel.From != "USER" || rel.To != "ORDER" {
		t.Errorf("Expected relationship from USER to ORDER, got %s to %s", rel.From, rel.To)
	}
	
	if rel.Type != OneToMany {
		t.Errorf("Expected OneToMany relationship, got %v", rel.Type)
	}
}

func TestDetectDiagramType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected DiagramType
	}{
		{
			name:     "ER diagram",
			input:    "erDiagram\n    USER {}",
			expected: ERDiagramType,
		},
		{
			name:     "Sequence diagram",
			input:    "sequenceDiagram\n    A->B: Hello",
			expected: SequenceDiagramType,
		},
		{
			name:     "Default to sequence",
			input:    "A->B: Hello",
			expected: SequenceDiagramType,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectDiagramType(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseRelationshipSymbol(t *testing.T) {
	tests := []struct {
		name         string
		symbol       string
		expectedType RelationshipType
		expectedFrom string
		expectedTo   string
	}{
		{
			name:         "one to one",
			symbol:       "||--||",
			expectedType: OneToOne,
			expectedFrom: "1",
			expectedTo:   "1",
		},
		{
			name:         "one to many",
			symbol:       "||--o{",
			expectedType: OneToMany,
			expectedFrom: "1",
			expectedTo:   "M",
		},
		{
			name:         "many to one",
			symbol:       "}o--||",
			expectedType: ManyToOne,
			expectedFrom: "M",
			expectedTo:   "1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relType, fromCard, toCard := parseRelationshipSymbol(tt.symbol)
			
			if relType != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, relType)
			}
			
			if fromCard != tt.expectedFrom {
				t.Errorf("Expected from cardinality %s, got %s", tt.expectedFrom, fromCard)
			}
			
			if toCard != tt.expectedTo {
				t.Errorf("Expected to cardinality %s, got %s", tt.expectedTo, toCard)
			}
		})
	}
}