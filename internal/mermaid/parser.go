package mermaid

import (
	"bufio"
	"regexp"
	"strings"
)

type DiagramType int

const (
	SequenceDiagramType DiagramType = iota
	ERDiagramType
)

type Diagram interface {
	GetType() DiagramType
}

type SequenceDiagram struct {
	Participants []Participant
	Messages     []Message
}

func (sd *SequenceDiagram) GetType() DiagramType {
	return SequenceDiagramType
}

type Participant struct {
	Name  string
	Alias string
}

type Message struct {
	From    string
	To      string
	Text    string
	Type    MessageType
	Activate bool
	Deactivate bool
}

type MessageType int

const (
	SolidArrow MessageType = iota
	DashedArrow
	SolidArrowWithX
	DashedArrowWithX
	Note
)

type ERDiagram struct {
	Entities      []Entity
	Relationships []Relationship
}

func (erd *ERDiagram) GetType() DiagramType {
	return ERDiagramType
}

type Entity struct {
	Name       string
	Attributes []Attribute
}

type Attribute struct {
	Name       string
	Type       string
	IsPK       bool
	IsFK       bool
	IsUnique   bool
	IsNotNull  bool
}

type Relationship struct {
	From         string
	To           string
	Type         RelationshipType
	FromCardinality string
	ToCardinality   string
	Label        string
}

type RelationshipType int

const (
	OneToOne RelationshipType = iota
	OneToMany
	ManyToOne
	ManyToMany
	Identifying
	NonIdentifying
)

func ParseDiagram(input string) (Diagram, error) {
	diagramType := DetectDiagramType(input)
	
	switch diagramType {
	case SequenceDiagramType:
		return ParseSequenceDiagram(input)
	case ERDiagramType:
		return ParseERDiagram(input)
	default:
		return ParseSequenceDiagram(input) // Default to sequence diagram
	}
}

func DetectDiagramType(input string) DiagramType {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "erDiagram") {
			return ERDiagramType
		}
		if strings.HasPrefix(line, "sequenceDiagram") {
			return SequenceDiagramType
		}
	}
	return SequenceDiagramType // Default
}

func ParseERDiagram(input string) (*ERDiagram, error) {
	diagram := &ERDiagram{
		Entities:      make([]Entity, 0),
		Relationships: make([]Relationship, 0),
	}
	
	scanner := bufio.NewScanner(strings.NewReader(input))
	entityMap := make(map[string]*Entity)
	var currentEntity *Entity
	inEntity := false
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if shouldSkipLine(line) {
			continue
		}
		
		if line == "}" && inEntity {
			if currentEntity != nil {
				diagram.Entities = append(diagram.Entities, *currentEntity)
				entityMap[currentEntity.Name] = currentEntity
			}
			inEntity = false
			currentEntity = nil
			continue
		}
		
		if entity := parseEntityStart(line); entity != nil {
			currentEntity = entity
			inEntity = true
			
			// Check if the entity definition ends on the same line (e.g., "USER {}")
			if strings.HasSuffix(line, "{}") {
				diagram.Entities = append(diagram.Entities, *currentEntity)
				entityMap[currentEntity.Name] = currentEntity
				inEntity = false
				currentEntity = nil
			}
			continue
		}
		
		if inEntity && currentEntity != nil {
			if attr := parseAttribute(line); attr != nil {
				currentEntity.Attributes = append(currentEntity.Attributes, *attr)
			}
			continue
		}
		
		if rel := parseRelationship(line); rel != nil {
			diagram.Relationships = append(diagram.Relationships, *rel)
		}
	}
	
	if inEntity && currentEntity != nil {
		diagram.Entities = append(diagram.Entities, *currentEntity)
	}
	
	return diagram, scanner.Err()
}

func shouldSkipLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "%%") || strings.HasPrefix(line, "erDiagram")
}

func parseEntityStart(line string) *Entity {
	entityRegex := regexp.MustCompile(`^\s*(\w+)\s*\{`)
	if matches := entityRegex.FindStringSubmatch(line); matches != nil {
		return &Entity{
			Name:       matches[1],
			Attributes: make([]Attribute, 0),
		}
	}
	return nil
}

func parseAttribute(line string) *Attribute {
	attributeRegex := regexp.MustCompile(`^\s*(\w+)\s+(\w+)(?:\s+(PK|FK|UK|NOT NULL))*`)
	if matches := attributeRegex.FindStringSubmatch(line); matches != nil {
		attr := &Attribute{
			Name: matches[2],
			Type: matches[1],
		}
		
		if len(matches) > 3 {
			constraints := strings.ToUpper(matches[3])
			attr.IsPK = strings.Contains(constraints, "PK")
			attr.IsFK = strings.Contains(constraints, "FK")
			attr.IsUnique = strings.Contains(constraints, "UK")
			attr.IsNotNull = strings.Contains(constraints, "NOT NULL")
		}
		
		return attr
	}
	return nil
}

func parseRelationship(line string) *Relationship {
	if !strings.Contains(line, "--") || !strings.Contains(line, ":") {
		return nil
	}
	
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return nil
	}
	
	leftPart := strings.TrimSpace(parts[0])
	label := strings.TrimSpace(parts[1])
	
	words := strings.Fields(leftPart)
	if len(words) < 3 {
		return nil
	}
	
	from := words[0]
	to := words[len(words)-1]
	relationSymbol := strings.Join(words[1:len(words)-1], " ")
	
	relType, fromCard, toCard := parseRelationshipSymbol(relationSymbol)
	
	return &Relationship{
		From:            from,
		To:              to,
		Type:            relType,
		FromCardinality: fromCard,
		ToCardinality:   toCard,
		Label:           label,
	}
}

func shouldSkipSequenceLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "%%") || strings.HasPrefix(line, "sequenceDiagram")
}

func parseParticipant(line string) *Participant {
	participantRegex := regexp.MustCompile(`^\s*participant\s+(\w+)(?:\s+as\s+(.+))?`)
	if matches := participantRegex.FindStringSubmatch(line); matches != nil {
		name := matches[1]
		alias := name
		if len(matches) > 2 && matches[2] != "" {
			alias = matches[2]
		}
		return &Participant{Name: name, Alias: alias}
	}
	return nil
}

func parseMessage(line string) *Message {
	messageRegex := regexp.MustCompile(`^\s*(\w+)\s*(->|-->|->>|-->>)\s*(\w+)\s*:\s*(.*)`)
	if matches := messageRegex.FindStringSubmatch(line); matches != nil {
		from := matches[1]
		arrow := matches[2]
		to := matches[3]
		text := matches[4]
		
		msgType := getMessageType(arrow)
		
		return &Message{
			From: from,
			To:   to,
			Text: text,
			Type: msgType,
		}
	}
	return nil
}

func getMessageType(arrow string) MessageType {
	switch arrow {
	case "->":
		return SolidArrow
	case "-->":
		return DashedArrow
	case "->>":
		return SolidArrowWithX
	case "-->>":
		return DashedArrowWithX
	default:
		return SolidArrow
	}
}

func ensureParticipantsExist(diagram *SequenceDiagram, participantMap map[string]bool, participants ...string) {
	for _, p := range participants {
		if !participantMap[p] {
			diagram.Participants = append(diagram.Participants, Participant{
				Name:  p,
				Alias: p,
			})
			participantMap[p] = true
		}
	}
}

func parseActivation(line string) {
	activateRegex := regexp.MustCompile(`^\s*activate\s+(\w+)`)
	deactivateRegex := regexp.MustCompile(`^\s*deactivate\s+(\w+)`)
	
	if activateRegex.MatchString(line) {
		// Handle activation - for now, skip
		return
	}
	
	if deactivateRegex.MatchString(line) {
		// Handle deactivation - for now, skip
		return
	}
}

func parseRelationshipSymbol(symbol string) (RelationshipType, string, string) {
	// Simplified relationship parsing
	// ||--|| : one to one
	// ||--o{ : one to many
	// }o--|| : many to one
	// }o--o{ : many to many
	
	fromCard := "1"
	toCard := "1"
	relType := OneToOne
	
	// Look for curly braces indicating "many"
	if strings.Contains(symbol, "{") {
		toCard = "M"
	}
	if strings.Contains(symbol, "}") {
		fromCard = "M"
	}
	
	// Determine relationship type based on cardinalities
	if fromCard == "1" && toCard == "1" {
		relType = OneToOne
	} else if fromCard == "1" && toCard == "M" {
		relType = OneToMany
	} else if fromCard == "M" && toCard == "1" {
		relType = ManyToOne
	} else if fromCard == "M" && toCard == "M" {
		relType = ManyToMany
	}
	
	return relType, fromCard, toCard
}

func ParseSequenceDiagram(input string) (*SequenceDiagram, error) {
	diagram := &SequenceDiagram{
		Participants: make([]Participant, 0),
		Messages:     make([]Message, 0),
	}
	
	scanner := bufio.NewScanner(strings.NewReader(input))
	participantMap := make(map[string]bool)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if shouldSkipSequenceLine(line) {
			continue
		}
		
		if participant := parseParticipant(line); participant != nil {
			if !participantMap[participant.Name] {
				diagram.Participants = append(diagram.Participants, *participant)
				participantMap[participant.Name] = true
			}
			continue
		}
		
		if message := parseMessage(line); message != nil {
			ensureParticipantsExist(diagram, participantMap, message.From, message.To)
			diagram.Messages = append(diagram.Messages, *message)
			continue
		}
		
		parseActivation(line)
	}
	
	return diagram, scanner.Err()
}