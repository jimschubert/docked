package types

import (
	"fmt"
	"strconv"
	"strings"
)


// intMax returns the max of two ints
func intMax(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// intMin returns the min of two ints
func intMin(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Protocol references the defined protocol (tcp, udp, or empty)
type Protocol string

// ExposeList abstracts on a slice of Expose instances
type ExposeList []Expose

// Expose is a typed representation of the docker EXPOSE command
type Expose struct {
	Original string `json:"original,omitempty"`
	PortRange PortRange `json:"port_range"`
	Protocol Protocol `json:"protocol,omitempty"`
	Description *string `json:"description,omitempty"`
}

// PortRange defines the start and end points for exposed ports.
// If a single port is exposed, Start and End will be the same value.
type PortRange struct {
	Start int
	End int
}

// Parse an expose definition string into a typed Expose instance, error if port definition is not valid
//goland:noinspection GoBoolExpressions
func (e *Expose) Parse(def string) error {
	parts := strings.FieldsFunc(def, func(r rune) bool {
		return r == '/'
	})
	portRange := PortRange{}
	portRange.Of(parts[0])
	if !portRange.IsValid() {
		return fmt.Errorf("invalid port range: start=(%d) end=(%d)", portRange.Start, portRange.End)
	}

	if len(parts) > 1 {
		userProtocol := parts[1]
		protocol := strings.ToLower(userProtocol)
		if protocol != "tcp" && protocol != "udp" && protocol != "" {
			return fmt.Errorf("invalid protocol used in EXPOSE command: %s", parts[1])
		}
		e.Protocol = Protocol(protocol)
	}

	e.Original = def
	e.PortRange = portRange

	return nil
}

// Of returns the valid PortRange representation of the pass string input
func (p *PortRange) Of(input string) *PortRange {
	ports := strings.FieldsFunc(input, func(r rune) bool {
		return r == '-'
	})
	for i, port := range ports {
		parsed, err := strconv.Atoi(port)
		if err == nil {
			if i == 0 {
				p.Start = parsed
			} else {
				p.End = parsed
			}
		}
	}
	if len(ports) == 1 {
		p.End = p.Start
	}

	return p
}

// IsValid determines if the defined ports are within the valid port range (a short)
func (p *PortRange) IsValid() bool {
	return 0 < p.Start && p.Start <= 65535 && 0 < p.End && p.End <= 65535
}

// Intersects determines if two ranges intersect (overlap)
//
//	 pA           pB
//	  ┌────────────┐
//	                  oA          oB
//	                  ┌────────────┐
//
//	  max(pA,oA) > min(pB, oB)
//	  * No intersection
//
//	 pA           pB
//	  ┌────────────┐
//	          oA          oB
//	          ┌────────────┐
//
//	  max(pA,oA) <= min(pB, oB)
//	  * These intersect
func (p *PortRange) Intersects(other PortRange) bool {
	low := intMax(p.Start, other.Start)
	high := intMin(p.End, other.End)
	return low <= high
}

// ParseExposeList converts the original input string into a list of Expose instances, or an error if parsing fails
func ParseExposeList(original string) (ExposeList, error)  {
	exposed := ExposeList{}
	defs := strings.Fields(original)
	for _, def := range defs {
		expose := Expose{}
		err := expose.Parse(def)
		if err != nil {
			return nil, err
		}
		exposed = append(exposed, expose)
	}

	return exposed, nil
}

