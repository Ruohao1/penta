package ports

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Ruohao1/penta/internal/model"
)

func Resolve(portExpr []string) ([]model.Port, error) {
	partsExpr := strings.TrimSpace(strings.Join(portExpr, ","))
	if partsExpr == "" {
		return defaultPorts(), nil
	}

	parts := strings.Split(partsExpr, ",")

	ports := make([]model.Port, 0, 64)

	for _, expr := range parts {
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}

		// aliases
		if expr == "all" || expr == "-" {
			return allPorts(), nil // shortcut: ignore other tokens
		}

		if strings.Contains(expr, "-") {
			portRange, err := expandRange(expr)
			if err != nil {
				return nil, fmt.Errorf("parse %q as port range: %w", expr, err)
			}
			ports = append(ports, portRange...)
			continue
		}

		p, err := strconv.Atoi(expr)
		if err != nil {
			return nil, fmt.Errorf("parse %q as port: %w", expr, err)
		}
		if p < 1 || p > 65535 {
			return nil, fmt.Errorf("invalid port %d", p)
		}
		ports = append(ports, model.Port{Number: p})
	}

	if len(ports) == 0 {
		return defaultPorts(), nil
	}

	return dedupeSort(ports), nil
}

func expandRange(expr string) ([]model.Port, error) {
	// enforce exactly "start-end"
	if strings.Count(expr, "-") != 1 {
		return nil, fmt.Errorf("invalid port range %q", expr)
	}
	a, b, _ := strings.Cut(expr, "-")
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if a == "" || b == "" {
		return nil, fmt.Errorf("invalid port range %q", expr)
	}

	start, err := strconv.Atoi(a)
	if err != nil {
		return nil, fmt.Errorf("invalid port range %q", expr)
	}
	end, err := strconv.Atoi(b)
	if err != nil {
		return nil, fmt.Errorf("invalid port range %q", expr)
	}
	if start < 1 || start > 65535 || end < 1 || end > 65535 || start > end {
		return nil, fmt.Errorf("invalid port range %q", expr)
	}

	out := make([]model.Port, 0, end-start+1)
	for i := start; i <= end; i++ {
		out = append(out, model.Port{Number: i})
	}
	return out, nil
}

func allPorts() []model.Port {
	out := make([]model.Port, 0, 65535)
	for p := 1; p <= 65535; p++ {
		out = append(out, model.Port{Number: p})
	}
	return out
}

func dedupeSort(in []model.Port) []model.Port {
	seen := make(map[int]struct{}, len(in))
	out := make([]model.Port, 0, len(in))

	for _, p := range in {
		if _, ok := seen[p.Number]; ok {
			continue
		}
		seen[p.Number] = struct{}{}
		out = append(out, model.Port{Number: p.Number})
	}

	slices.SortFunc(out, func(a, b model.Port) int { return a.Number - b.Number })
	return out
}

func defaultPorts() []model.Port {
	return []model.Port{{Number: 22}, {Number: 80}, {Number: 443}}
}
