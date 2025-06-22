package drawio

import (
	"encoding/xml"
	"fmt"
	"mermaid2drawio/internal/mermaid"
)

// Layout constants for entity diagrams
const (
	EntityWidth      = 200.0
	EntityHeight     = 30.0
	AttributeHeight  = 20.0
	EntitySpacing    = 300.0
	StartX           = 50.0
	StartY           = 50.0
	EntitiesPerRow   = 3
)

// Layout constants for sequence diagrams
const (
	ParticipantWidth   = 100.0
	ParticipantHeight  = 50.0
	ParticipantSpacing = 200.0
	ParticipantY       = 50.0
)

// Draw.io model defaults
const (
	DefaultDx         = 1234
	DefaultDy         = 732
	DefaultGrid       = 1
	DefaultGridSize   = 10
	DefaultGuides     = 1
	DefaultTooltips   = 1
	DefaultConnect    = 1
	DefaultArrows     = 1
	DefaultFold       = 1
	DefaultPage       = 1
	DefaultPageScale  = 1.0
	DefaultPageWidth  = 827
	DefaultPageHeight = 1169
	DefaultMath       = 0
	DefaultShadow     = 0
)

type MxGraphModel struct {
	XMLName xml.Name `xml:"mxGraphModel"`
	Dx      int      `xml:"dx,attr"`
	Dy      int      `xml:"dy,attr"`
	Grid    int      `xml:"grid,attr"`
	GridSize int     `xml:"gridSize,attr"`
	Guides  int      `xml:"guides,attr"`
	Tooltips int     `xml:"tooltips,attr"`
	Connect int      `xml:"connect,attr"`
	Arrows  int      `xml:"arrows,attr"`
	Fold    int      `xml:"fold,attr"`
	Page    int      `xml:"page,attr"`
	PageScale float64 `xml:"pageScale,attr"`
	PageWidth int     `xml:"pageWidth,attr"`
	PageHeight int    `xml:"pageHeight,attr"`
	Math    int      `xml:"math,attr"`
	Shadow  int      `xml:"shadow,attr"`
	Root    *MxRoot  `xml:"root"`
}

type MxRoot struct {
	MxCells []MxCell `xml:"mxCell"`
}

type MxCell struct {
	ID       string      `xml:"id,attr"`
	Value    string      `xml:"value,attr,omitempty"`
	Style    string      `xml:"style,attr,omitempty"`
	Vertex   string      `xml:"vertex,attr,omitempty"`
	Edge     string      `xml:"edge,attr,omitempty"`
	Parent   string      `xml:"parent,attr,omitempty"`
	Source   string      `xml:"source,attr,omitempty"`
	Target   string      `xml:"target,attr,omitempty"`
	Geometry *MxGeometry `xml:"mxGeometry,omitempty"`
}

type MxGeometry struct {
	X      *float64 `xml:"x,attr,omitempty"`
	Y      *float64 `xml:"y,attr,omitempty"`
	Width  *float64 `xml:"width,attr,omitempty"`
	Height *float64 `xml:"height,attr,omitempty"`
	As     string   `xml:"as,attr"`
}

func createBaseModel() *MxGraphModel {
	return &MxGraphModel{
		Dx:         DefaultDx,
		Dy:         DefaultDy,
		Grid:       DefaultGrid,
		GridSize:   DefaultGridSize,
		Guides:     DefaultGuides,
		Tooltips:   DefaultTooltips,
		Connect:    DefaultConnect,
		Arrows:     DefaultArrows,
		Fold:       DefaultFold,
		Page:       DefaultPage,
		PageScale:  DefaultPageScale,
		PageWidth:  DefaultPageWidth,
		PageHeight: DefaultPageHeight,
		Math:       DefaultMath,
		Shadow:     DefaultShadow,
		Root:       &MxRoot{},
	}
}

func generateXMLOutput(model *MxGraphModel) (string, error) {
	output, err := xml.MarshalIndent(model, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}

func createDefaultCells() []MxCell {
	return []MxCell{
		{ID: "0"},
		{ID: "1", Parent: "0"},
	}
}

func GenerateDrawIOXML(diagram mermaid.Diagram) (string, error) {
	switch d := diagram.(type) {
	case *mermaid.SequenceDiagram:
		return GenerateSequenceDrawIOXML(d)
	case *mermaid.ERDiagram:
		return GenerateERDrawIOXML(d)
	default:
		return "", fmt.Errorf("unsupported diagram type")
	}
}

func GenerateERDrawIOXML(diagram *mermaid.ERDiagram) (string, error) {
	model := createBaseModel()

	cells := createDefaultCells()
	cellID := 2
	entityCells := make(map[string]string)
	
	// Create entity tables
	currentRow := 0
	currentCol := 0
	
	for _, entity := range diagram.Entities {
		x := StartX + float64(currentCol)*EntitySpacing
		y := StartY + float64(currentRow)*EntitySpacing
		
		totalHeight := EntityHeight + float64(len(entity.Attributes))*AttributeHeight
		
		// Create entity header
		headerID := fmt.Sprintf("entity_header_%d", cellID)
		entityCells[entity.Name] = headerID
		
		entityWidth := EntityWidth
		headerCell := MxCell{
			ID:     headerID,
			Value:  entity.Name,
			Style:  "swimlane;fontStyle=1;align=center;verticalAlign=middle;childLayout=stackLayout;horizontal=1;startSize=30;horizontalStack=0;resizeParent=1;resizeParentMax=0;resizeLast=0;collapsible=0;marginBottom=0;whiteSpace=wrap;html=1;",
			Vertex: "1",
			Parent: "1",
			Geometry: &MxGeometry{
				X:      &x,
				Y:      &y,
				Width:  &entityWidth,
				Height: &totalHeight,
				As:     "geometry",
			},
		}
		cells = append(cells, headerCell)
		cellID++
		
		// Create attributes
		for i, attr := range entity.Attributes {
			attrY := y + EntityHeight + float64(i)*AttributeHeight
			attrID := fmt.Sprintf("attr_%d", cellID)
			
			// Format attribute text with constraints
			attrText := fmt.Sprintf("%s: %s", attr.Name, attr.Type)
			if attr.IsPK {
				attrText = "🔑 " + attrText
			}
			if attr.IsFK {
				attrText = "🔗 " + attrText
			}
			if attr.IsUnique {
				attrText += " (UK)"
			}
			if attr.IsNotNull {
				attrText += " (NN)"
			}
			
			entityWidth := EntityWidth
			attributeHeight := AttributeHeight
			attrCell := MxCell{
				ID:     attrID,
				Value:  attrText,
				Style:  "text;strokeColor=none;fillColor=none;align=left;verticalAlign=middle;spacingLeft=4;spacingRight=4;overflow=hidden;points=[[0,0.5],[1,0.5]];portConstraint=eastwest;rotatable=0;whiteSpace=wrap;html=1;",
				Vertex: "1",
				Parent: headerID,
				Geometry: &MxGeometry{
					Y:      &attrY,
					Width:  &entityWidth,
					Height: &attributeHeight,
					As:     "geometry",
				},
			}
			cells = append(cells, attrCell)
			cellID++
		}
		
		// Update position for next entity
		currentCol++
		if currentCol >= EntitiesPerRow {
			currentCol = 0
			currentRow++
		}
	}
	
	// Create relationships
	for _, relationship := range diagram.Relationships {
		fromID := entityCells[relationship.From]
		toID := entityCells[relationship.To]
		
		if fromID == "" || toID == "" {
			continue // Skip if entity not found
		}
		
		// Determine relationship style
		style := "endArrow=none;html=1;rounded=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;exitX=1;exitY=0.5;exitDx=0;exitDy=0;"
		
		// Add cardinality indicators
		switch relationship.Type {
		case mermaid.OneToOne:
			style += "startArrow=ERone;endArrow=ERone;"
		case mermaid.OneToMany:
			style += "startArrow=ERone;endArrow=ERmany;"
		case mermaid.ManyToOne:
			style += "startArrow=ERmany;endArrow=ERone;"
		case mermaid.ManyToMany:
			style += "startArrow=ERmany;endArrow=ERmany;"
		}
		
		relationshipCell := MxCell{
			ID:     fmt.Sprintf("relationship_%d", cellID),
			Value:  relationship.Label,
			Style:  style,
			Edge:   "1",
			Parent: "1",
			Source: fromID,
			Target: toID,
			Geometry: &MxGeometry{
				As: "geometry",
			},
		}
		cells = append(cells, relationshipCell)
		cellID++
	}
	
	model.Root.MxCells = cells
	
	// Generate XML
	output, err := xml.MarshalIndent(model, "", "  ")
	if err != nil {
		return "", err
	}
	
	return xml.Header + string(output), nil
}

func GenerateSequenceDrawIOXML(diagram *mermaid.SequenceDiagram) (string, error) {
	model := createBaseModel()

	cells := createDefaultCells()
	cellID := 2
	participantCells := make(map[string]string)
	
	// Create participant rectangles (actors)
	
	for i, participant := range diagram.Participants {
		x := StartX + float64(i)*ParticipantSpacing
		id := fmt.Sprintf("participant_%d", cellID)
		participantCells[participant.Name] = id
		
		participantY := ParticipantY
		participantWidth := ParticipantWidth
		participantHeight := ParticipantHeight
		cell := MxCell{
			ID:     id,
			Value:  participant.Alias,
			Style:  "rounded=0;whiteSpace=wrap;html=1;",
			Vertex: "1",
			Parent: "1",
			Geometry: &MxGeometry{
				X:      &x,
				Y:      &participantY,
				Width:  &participantWidth,
				Height: &participantHeight,
				As:     "geometry",
			},
		}
		cells = append(cells, cell)
		cellID++
		
		// Add lifeline (vertical line)
		lifelineCell := MxCell{
			ID:     fmt.Sprintf("lifeline_%d", cellID),
			Style:  "endArrow=none;dashed=1;html=1;",
			Edge:   "1",
			Parent: "1",
			Geometry: &MxGeometry{
				As: "geometry",
			},
		}
		cells = append(cells, lifelineCell)
		cellID++
	}
	
	// Create message arrows
	for _, message := range diagram.Messages {
		fromID := participantCells[message.From]
		toID := participantCells[message.To]
		
		if fromID == "" || toID == "" {
			continue // Skip if participant not found
		}
		
		// Determine arrow style based on message type
		var style string
		switch message.Type {
		case mermaid.SolidArrow:
			style = "endArrow=classic;html=1;"
		case mermaid.DashedArrow:
			style = "endArrow=classic;html=1;dashed=1;"
		case mermaid.SolidArrowWithX:
			style = "endArrow=block;endFill=1;html=1;"
		case mermaid.DashedArrowWithX:
			style = "endArrow=block;endFill=1;html=1;dashed=1;"
		default:
			style = "endArrow=classic;html=1;"
		}
		
		messageCell := MxCell{
			ID:     fmt.Sprintf("message_%d", cellID),
			Value:  message.Text,
			Style:  style,
			Edge:   "1",
			Parent: "1",
			Source: fromID,
			Target: toID,
			Geometry: &MxGeometry{
				As: "geometry",
			},
		}
		cells = append(cells, messageCell)
		cellID++
	}
	
	model.Root.MxCells = cells
	return generateXMLOutput(model)
}