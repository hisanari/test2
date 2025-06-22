package mermaid

import (
	"testing"
)

func TestParseSequenceDiagram(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedParticipants int
		expectedMessages     int
		wantErr           bool
	}{
		{
			name:              "simple diagram",
			input:             "sequenceDiagram\n    A->B: Hello",
			expectedParticipants: 2,
			expectedMessages:     1,
			wantErr:           false,
		},
		{
			name:              "with participants",
			input:             "sequenceDiagram\n    participant A as Alice\n    participant B as Bob\n    A->B: Hello",
			expectedParticipants: 2,
			expectedMessages:     1,
			wantErr:           false,
		},
		{
			name:              "multiple messages",
			input:             "sequenceDiagram\n    A->B: Hello\n    B-->A: Hi\n    A->>B: How are you?",
			expectedParticipants: 2,
			expectedMessages:     3,
			wantErr:           false,
		},
		{
			name:              "empty input",
			input:             "",
			expectedParticipants: 0,
			expectedMessages:     0,
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := ParseSequenceDiagram(tt.input)
			
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if len(diagram.Participants) != tt.expectedParticipants {
				t.Errorf("Expected %d participants, got %d", tt.expectedParticipants, len(diagram.Participants))
			}
			
			if len(diagram.Messages) != tt.expectedMessages {
				t.Errorf("Expected %d messages, got %d", tt.expectedMessages, len(diagram.Messages))
			}
		})
	}
}

func TestMessageTypes(t *testing.T) {
	input := `sequenceDiagram
    A->B: solid arrow
    A-->B: dashed arrow
    A->>B: solid arrow with x
    A-->>B: dashed arrow with x`
	
	diagram, err := ParseSequenceDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	expectedTypes := []MessageType{
		SolidArrow,
		DashedArrow,
		SolidArrowWithX,
		DashedArrowWithX,
	}
	
	if len(diagram.Messages) != len(expectedTypes) {
		t.Fatalf("Expected %d messages, got %d", len(expectedTypes), len(diagram.Messages))
	}
	
	for i, expectedType := range expectedTypes {
		if diagram.Messages[i].Type != expectedType {
			t.Errorf("Message %d: expected type %v, got %v", i, expectedType, diagram.Messages[i].Type)
		}
	}
}

func TestGetType(t *testing.T) {
	seqDiagram := &SequenceDiagram{}
	if seqDiagram.GetType() != SequenceDiagramType {
		t.Error("SequenceDiagram should return SequenceDiagramType")
	}

	erDiagram := &ERDiagram{}
	if erDiagram.GetType() != ERDiagramType {
		t.Error("ERDiagram should return ERDiagramType")
	}
}

func TestParseDiagram(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected DiagramType
		wantErr  bool
	}{
		{
			name:     "sequence diagram",
			input:    "sequenceDiagram\n    A->B: Hello",
			expected: SequenceDiagramType,
			wantErr:  false,
		},
		{
			name:     "ER diagram",
			input:    "erDiagram\n    USER {}",
			expected: ERDiagramType,
			wantErr:  false,
		},
		{
			name:     "default to sequence",
			input:    "unknown diagram\n    A->B: Hello",
			expected: SequenceDiagramType,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := ParseDiagram(tt.input)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if diagram.GetType() != tt.expected {
				t.Errorf("Expected type %v, got %v", tt.expected, diagram.GetType())
			}
		})
	}
}

func TestParseERDiagramWithAttributes(t *testing.T) {
	input := `erDiagram
    USER {
        int id PK
        string name NOT NULL
        string email UK
        int role_id FK
    }`

	diagram, err := ParseERDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(diagram.Entities) != 1 {
		t.Fatalf("Expected 1 entity, got %d", len(diagram.Entities))
	}

	entity := diagram.Entities[0]
	if entity.Name != "USER" {
		t.Errorf("Expected entity name USER, got %s", entity.Name)
	}

	if len(entity.Attributes) != 4 {
		t.Fatalf("Expected 4 attributes, got %d", len(entity.Attributes))
	}

	// Check specific constraints
	pkAttr := entity.Attributes[0]
	if !pkAttr.IsPK {
		t.Error("id attribute should be marked as PK")
	}

	fkAttr := entity.Attributes[3]
	if !fkAttr.IsFK {
		t.Error("role_id attribute should be marked as FK")
	}
}

func TestParseERDiagramComplex(t *testing.T) {
	input := `erDiagram
    USER {
        int id PK
    }
    ORDER {
        int id PK
        int user_id FK
    }
    PRODUCT {
        int id PK
    }
    USER ||--o{ ORDER : places
    ORDER }o--|| PRODUCT : contains`

	diagram, err := ParseERDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(diagram.Entities) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(diagram.Entities))
	}

	if len(diagram.Relationships) != 2 {
		t.Errorf("Expected 2 relationships, got %d", len(diagram.Relationships))
	}
}

func TestParseActivation(t *testing.T) {
	// Test activation/deactivation parsing (currently skipped)
	input := `sequenceDiagram
    A->B: Hello
    activate B
    B->A: Hi
    deactivate B`

	diagram, err := ParseSequenceDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should still parse messages correctly despite activation commands
	if len(diagram.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(diagram.Messages))
	}
}

func TestParseMessageDefaultType(t *testing.T) {
	input := `sequenceDiagram
    A->B: solid
    A-->B: dashed
    A->>B: solid_x
    A-->>B: dashed_x
    A-?B: unknown`

	diagram, err := ParseSequenceDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTypes := []MessageType{
		SolidArrow,
		DashedArrow,
		SolidArrowWithX,
		DashedArrowWithX,
		SolidArrow, // unknown type defaults to SolidArrow
	}

	for i, expectedType := range expectedTypes {
		if i < len(diagram.Messages) && diagram.Messages[i].Type != expectedType {
			t.Errorf("Message %d: expected type %v, got %v", i, expectedType, diagram.Messages[i].Type)
		}
	}
}

func TestParseRelationshipEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"no colon", "USER ||--o{ ORDER", false},
		{"no relationship symbol", "USER ORDER : test", false},
		{"too few words", "USER : test", false},
		{"valid relationship", "USER ||--o{ ORDER : places", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := "erDiagram\n    USER {}\n    ORDER {}\n    " + tt.input
			diagram, err := ParseERDiagram(input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			hasRelationship := len(diagram.Relationships) > 0
			if tt.valid != hasRelationship {
				t.Errorf("Expected valid=%v, got relationship count=%d", tt.valid, len(diagram.Relationships))
			}
		})
	}
}

func TestParseRelationshipSymbolEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		symbol       string
		expectedType RelationshipType
	}{
		{"many to many", "}o--o{", ManyToMany},
		{"identifying", "||--||", OneToOne},
		{"non-identifying", "||..||", OneToOne},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relType, _, _ := parseRelationshipSymbol(tt.symbol)
			if relType != tt.expectedType {
				t.Errorf("Expected %v, got %v", tt.expectedType, relType)
			}
		})
	}
}

func TestParseDiagramDefaultCase(t *testing.T) {
	// Test the default case in ParseDiagram
	input := "unknownDiagram\n    unknown content"
	
	diagram, err := ParseDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should default to sequence diagram
	if diagram.GetType() != SequenceDiagramType {
		t.Errorf("Expected SequenceDiagramType, got %v", diagram.GetType())
	}
}

func TestParseAttributeEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectMatch bool
	}{
		{"valid attribute", "int id PK", true},
		{"no constraint", "string name", true},
		{"invalid format", "invalid", false},
		{"only type", "int", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := parseAttribute(tt.input)
			hasMatch := attr != nil
			if hasMatch != tt.expectMatch {
				t.Errorf("Expected match=%v for input %q, got %v", tt.expectMatch, tt.input, hasMatch)
			}
		})
	}
}

func TestParseRelationshipInvalidCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no colon", "USER RELATIONSHIP ORDER"},
		{"no dash", "USER RELATIONSHIP ORDER : test"},
		{"only one part", "USER : test"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rel := parseRelationship(tt.input)
			if rel != nil {
				t.Errorf("Expected nil relationship for invalid input: %q", tt.input)
			}
		})
	}
}

func TestGetMessageTypeDefault(t *testing.T) {
	// Test default case in getMessageType
	msgType := getMessageType("unknown-arrow")
	if msgType != SolidArrow {
		t.Errorf("Expected SolidArrow for unknown type, got %v", msgType)
	}
}

func TestParseERDiagramEndInEntity(t *testing.T) {
	// Test case where input ends while still in entity
	input := `erDiagram
    USER {
        int id PK
        string name`
	
	diagram, err := ParseERDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should still create the entity despite incomplete structure
	if len(diagram.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(diagram.Entities))
	}
	
	if len(diagram.Entities[0].Attributes) != 2 {
		t.Errorf("Expected 2 attributes, got %d", len(diagram.Entities[0].Attributes))
	}
}

func TestParseSequenceDiagramActivationPaths(t *testing.T) {
	// Test activation parsing paths
	input := `sequenceDiagram
    A->B: Hello
    activate A
    deactivate A
    note right of A: This is a note`
	
	diagram, err := ParseSequenceDiagram(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should parse message correctly and ignore activation/note commands
	if len(diagram.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(diagram.Messages))
	}
}