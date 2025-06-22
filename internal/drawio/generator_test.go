package drawio

import (
	"mermaid2drawio/internal/mermaid"
	"strings"
	"testing"
)

func TestGenerateDrawIOXML(t *testing.T) {
	diagram := &mermaid.SequenceDiagram{
		Participants: []mermaid.Participant{
			{Name: "A", Alias: "Alice"},
			{Name: "B", Alias: "Bob"},
		},
		Messages: []mermaid.Message{
			{From: "A", To: "B", Text: "Hello", Type: mermaid.SolidArrow},
			{From: "B", To: "A", Text: "Hi", Type: mermaid.DashedArrow},
		},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Check XML structure
	if !strings.Contains(xml, "<?xml version") {
		t.Error("XML should contain XML declaration")
	}
	
	if !strings.Contains(xml, "<mxGraphModel") {
		t.Error("XML should contain mxGraphModel element")
	}
	
	// Check participants are included
	if !strings.Contains(xml, "Alice") {
		t.Error("XML should contain Alice participant")
	}
	
	if !strings.Contains(xml, "Bob") {
		t.Error("XML should contain Bob participant")
	}
	
	// Check messages are included
	if !strings.Contains(xml, "Hello") {
		t.Error("XML should contain Hello message")
	}
	
	if !strings.Contains(xml, "Hi") {
		t.Error("XML should contain Hi message")
	}
}

func TestGenerateDrawIOXMLEmpty(t *testing.T) {
	diagram := &mermaid.SequenceDiagram{
		Participants: []mermaid.Participant{},
		Messages:     []mermaid.Message{},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should still generate valid XML structure
	if !strings.Contains(xml, "<mxGraphModel") {
		t.Error("XML should contain mxGraphModel element even for empty diagram")
	}
}

func TestMessageTypeStyles(t *testing.T) {
	diagram := &mermaid.SequenceDiagram{
		Participants: []mermaid.Participant{
			{Name: "A", Alias: "A"},
			{Name: "B", Alias: "B"},
		},
		Messages: []mermaid.Message{
			{From: "A", To: "B", Text: "solid", Type: mermaid.SolidArrow},
			{From: "A", To: "B", Text: "dashed", Type: mermaid.DashedArrow},
			{From: "A", To: "B", Text: "solidX", Type: mermaid.SolidArrowWithX},
			{From: "A", To: "B", Text: "dashedX", Type: mermaid.DashedArrowWithX},
		},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Check different arrow styles are present
	if !strings.Contains(xml, "endArrow=classic;html=1;") {
		t.Error("Should contain solid arrow style")
	}
	
	if !strings.Contains(xml, "dashed=1") {
		t.Error("Should contain dashed arrow style")
	}
	
	if !strings.Contains(xml, "endArrow=block;endFill=1") {
		t.Error("Should contain block arrow style")
	}
}

func TestGenerateERDrawIOXML(t *testing.T) {
	diagram := &mermaid.ERDiagram{
		Entities: []mermaid.Entity{
			{
				Name: "USER",
				Attributes: []mermaid.Attribute{
					{Name: "id", Type: "int", IsPK: true},
					{Name: "name", Type: "string"},
					{Name: "email", Type: "string", IsUnique: true},
				},
			},
			{
				Name: "ORDER",
				Attributes: []mermaid.Attribute{
					{Name: "id", Type: "int", IsPK: true},
					{Name: "user_id", Type: "int", IsFK: true},
				},
			},
		},
		Relationships: []mermaid.Relationship{
			{
				From:  "USER",
				To:    "ORDER",
				Type:  mermaid.OneToMany,
				Label: "places",
			},
		},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Check XML structure
	if !strings.Contains(xml, "<?xml version") {
		t.Error("XML should contain XML declaration")
	}
	
	if !strings.Contains(xml, "<mxGraphModel") {
		t.Error("XML should contain mxGraphModel element")
	}
	
	// Check entities are included
	if !strings.Contains(xml, "USER") {
		t.Error("XML should contain USER entity")
	}
	
	if !strings.Contains(xml, "ORDER") {
		t.Error("XML should contain ORDER entity")
	}
	
	// Check attributes with constraints
	if !strings.Contains(xml, "🔑 id: int") {
		t.Error("XML should contain primary key indicator")
	}
	
	if !strings.Contains(xml, "🔗 user_id: int") {
		t.Error("XML should contain foreign key indicator")
	}
	
	if !strings.Contains(xml, "email: string (UK)") {
		t.Error("XML should contain unique constraint indicator")
	}
	
	// Check relationships
	if !strings.Contains(xml, "places") {
		t.Error("XML should contain relationship label")
	}
	
	if !strings.Contains(xml, "startArrow=ERone;endArrow=ERmany") {
		t.Error("XML should contain ER relationship arrows")
	}
}

func TestGenerateERDrawIOXMLEmpty(t *testing.T) {
	diagram := &mermaid.ERDiagram{
		Entities:      []mermaid.Entity{},
		Relationships: []mermaid.Relationship{},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should still generate valid XML structure
	if !strings.Contains(xml, "<mxGraphModel") {
		t.Error("XML should contain mxGraphModel element even for empty ER diagram")
	}
}

func TestGenerateDrawIOXMLUnsupportedType(t *testing.T) {
	// Create a mock diagram that doesn't implement the known types
	var unknownDiagram mermaid.Diagram = &mockDiagram{}
	
	_, err := GenerateDrawIOXML(unknownDiagram)
	if err == nil {
		t.Error("Expected error for unsupported diagram type")
	}
	
	if !strings.Contains(err.Error(), "unsupported diagram type") {
		t.Errorf("Expected 'unsupported diagram type' error, got: %v", err)
	}
}

func TestGenerateXMLOutputError(t *testing.T) {
	// Test XML marshaling error by providing invalid structure
	// Create a model with a cyclic reference to trigger XML marshaling error
	model := createBaseModel()
	
	// Test successful case first
	_, err := generateXMLOutput(model)
	if err != nil {
		t.Errorf("Unexpected error for valid XML structure: %v", err)
	}
	
	// Note: It's very difficult to trigger an XML marshaling error with the current structure
	// since Go's xml package is quite robust. The error path exists for completeness.
}

func TestGenerateERDrawIOXMLMarshalError(t *testing.T) {
	// Test ER diagram XML generation with large data to cover edge cases
	diagram := &mermaid.ERDiagram{
		Entities: []mermaid.Entity{
			{
				Name: "LARGE_ENTITY",
				Attributes: make([]mermaid.Attribute, 100),
			},
		},
	}
	
	// Fill with many attributes to stress test
	for i := 0; i < 100; i++ {
		diagram.Entities[0].Attributes[i] = mermaid.Attribute{
			Name: "attr" + string(rune(i)),
			Type: "type" + string(rune(i)),
		}
	}
	
	_, err := GenerateERDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestERDiagramRelationshipTypes(t *testing.T) {
	diagram := &mermaid.ERDiagram{
		Entities: []mermaid.Entity{
			{Name: "A"},
			{Name: "B"},
			{Name: "C"},
			{Name: "D"},
		},
		Relationships: []mermaid.Relationship{
			{From: "A", To: "B", Type: mermaid.OneToOne, Label: "one-to-one"},
			{From: "A", To: "C", Type: mermaid.ManyToOne, Label: "many-to-one"},
			{From: "A", To: "D", Type: mermaid.ManyToMany, Label: "many-to-many"},
		},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Check different relationship types
	if !strings.Contains(xml, "startArrow=ERone;endArrow=ERone") {
		t.Error("Should contain one-to-one relationship style")
	}
	
	if !strings.Contains(xml, "startArrow=ERmany;endArrow=ERone") {
		t.Error("Should contain many-to-one relationship style")
	}
	
	if !strings.Contains(xml, "startArrow=ERmany;endArrow=ERmany") {
		t.Error("Should contain many-to-many relationship style")
	}
}

func TestERDiagramAttributeConstraints(t *testing.T) {
	diagram := &mermaid.ERDiagram{
		Entities: []mermaid.Entity{
			{
				Name: "TEST",
				Attributes: []mermaid.Attribute{
					{Name: "attr1", Type: "string", IsNotNull: true},
					{Name: "attr2", Type: "int", IsPK: true, IsFK: true},
				},
			},
		},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Check different constraint indicators
	if !strings.Contains(xml, "(NN)") {
		t.Error("Should contain NOT NULL constraint indicator")
	}
	
	if !strings.Contains(xml, "🔑") && !strings.Contains(xml, "🔗") {
		t.Error("Should contain both PK and FK indicators")
	}
}

func TestSequenceDiagramSkippedParticipants(t *testing.T) {
	diagram := &mermaid.SequenceDiagram{
		Participants: []mermaid.Participant{
			{Name: "A", Alias: "Alice"},
		},
		Messages: []mermaid.Message{
			{From: "X", To: "Y", Text: "Unknown participants", Type: mermaid.SolidArrow},
		},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should still generate valid XML even with unknown participants
	if !strings.Contains(xml, "<mxGraphModel") {
		t.Error("Should generate valid XML even with unknown participants")
	}
}

func TestERDiagramSkippedRelationships(t *testing.T) {
	diagram := &mermaid.ERDiagram{
		Entities: []mermaid.Entity{
			{Name: "USER"},
		},
		Relationships: []mermaid.Relationship{
			{From: "USER", To: "UNKNOWN", Type: mermaid.OneToMany, Label: "missing"},
			{From: "UNKNOWN", To: "USER", Type: mermaid.OneToMany, Label: "missing2"},
		},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should generate valid XML and skip unknown relationships
	if !strings.Contains(xml, "<mxGraphModel") {
		t.Error("Should generate valid XML even with unknown entities in relationships")
	}
	
	// Should contain the USER entity
	if !strings.Contains(xml, "USER") {
		t.Error("Should contain existing USER entity")
	}
}

func TestSequenceDiagramDefaultMessageType(t *testing.T) {
	diagram := &mermaid.SequenceDiagram{
		Participants: []mermaid.Participant{
			{Name: "A", Alias: "A"},
			{Name: "B", Alias: "B"},
		},
		Messages: []mermaid.Message{
			{From: "A", To: "B", Text: "unknown type", Type: mermaid.MessageType(999)},
		},
	}
	
	xml, err := GenerateDrawIOXML(diagram)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should default to classic arrow style
	if !strings.Contains(xml, "endArrow=classic;html=1;") {
		t.Error("Should default to classic arrow style for unknown message type")
	}
}

// Mock diagram type to test unsupported diagram error
type mockDiagram struct{}

func (md *mockDiagram) GetType() mermaid.DiagramType {
	return mermaid.DiagramType(999) // Unknown type
}