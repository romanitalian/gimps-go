package primenet

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	// PrimeNetServerURL is the base URL for PrimeNet server
	PrimeNetServerURL = "http://v5.mersenne.org/v5server/"
	// APIVersion is the PrimeNet API version
	APIVersion = 0.95
)

// Client represents a PrimeNet HTTP client
type Client struct {
	baseURL    string
	httpClient *http.Client
	computerGUID string
}

// NewClient creates a new PrimeNet client
func NewClient(computerGUID string) *Client {
	return &Client{
		baseURL:     PrimeNetServerURL,
		httpClient:  &http.Client{},
		computerGUID: computerGUID,
	}
}

// UpdateComputerInfo sends update computer info request
func (c *Client) UpdateComputerInfo(req *UpdateComputerInfo) (*UpdateComputerInfo, error) {
	args := c.formatArgs(OpUpdateComputerInfo, req)
	resp, err := c.sendRequest(args)
	if err != nil {
		return nil, err
	}
	
	result := &UpdateComputerInfo{}
	if err := c.parseResponse(resp, OpUpdateComputerInfo, result); err != nil {
		return nil, err
	}
	
	return result, nil
}

// GetAssignment requests a new assignment from the server
func (c *Client) GetAssignment(req *GetAssignment) (*GetAssignment, error) {
	args := c.formatArgs(OpGetAssignment, req)
	resp, err := c.sendRequest(args)
	if err != nil {
		return nil, err
	}
	
	result := &GetAssignment{}
	if err := c.parseResponse(resp, OpGetAssignment, result); err != nil {
		return nil, err
	}
	
	return result, nil
}

// RegisterAssignment registers an assignment with the server
func (c *Client) RegisterAssignment(req *RegisterAssignment) error {
	args := c.formatArgs(OpRegisterAssignment, req)
	_, err := c.sendRequest(args)
	return err
}

// AssignmentProgress sends progress update for an assignment
func (c *Client) AssignmentProgress(req *AssignmentProgress) error {
	args := c.formatArgs(OpAssignmentProgress, req)
	_, err := c.sendRequest(args)
	return err
}

// AssignmentResult sends result of an assignment
func (c *Client) AssignmentResult(req *AssignmentResult) error {
	args := c.formatArgs(OpAssignmentResult, req)
	_, err := c.sendRequest(args)
	return err
}

// formatArgs formats HTTP arguments from a PrimeNet packet
func (c *Client) formatArgs(operation int, pkt any) string {
	var sb strings.Builder
	
	// Common header
	sb.WriteString(fmt.Sprintf("v=%.2f&px=GIMPS", APIVersion))
	
	// Format operation-specific arguments
	switch operation {
	case OpUpdateComputerInfo:
		req := pkt.(*UpdateComputerInfo)
		sb.WriteString("&t=uc")
		sb.WriteString("&g=" + url.QueryEscape(req.ComputerGUID))
		if req.HardwareGUID != "" {
			sb.WriteString("&hg=" + url.QueryEscape(req.HardwareGUID))
		}
		if req.WindowsGUID != "" {
			sb.WriteString("&wg=" + url.QueryEscape(req.WindowsGUID))
		}
		if req.Application != "" {
			sb.WriteString("&a=" + url.QueryEscape(req.Application))
		}
		if req.CPUModel != "" {
			sb.WriteString("&c=" + url.QueryEscape(req.CPUModel))
		}
		if req.CPUFeatures != "" {
			sb.WriteString("&f=" + url.QueryEscape(req.CPUFeatures))
		}
		sb.WriteString(fmt.Sprintf("&L1=%d&L2=%d&np=%d&hp=%d&m=%d&s=%d&h=%d&r=%d",
			req.L1CacheSize, req.L2CacheSize, req.NumCPUs,
			req.NumHyperthread, req.MemInstalled, req.CPUSpeed,
			req.HoursPerDay, req.RollingAverage))
		if req.L3CacheSize > 0 {
			sb.WriteString(fmt.Sprintf("&L3=%d", req.L3CacheSize))
		}
		if req.UserID != "" {
			sb.WriteString("&u=" + url.QueryEscape(req.UserID))
		}
		if req.ComputerName != "" {
			sb.WriteString("&cn=" + url.QueryEscape(req.ComputerName))
		}
		
	case OpGetAssignment:
		req := pkt.(*GetAssignment)
		sb.WriteString("&t=ga")
		sb.WriteString("&g=" + url.QueryEscape(req.ComputerGUID))
		sb.WriteString(fmt.Sprintf("&c=%d", req.CPUNum))
		if req.GetCertWork > 0 {
			sb.WriteString(fmt.Sprintf("&cert=%d", req.GetCertWork))
		}
		if req.TempDiskSpace > 0 {
			sb.WriteString(fmt.Sprintf("&tds=%f", req.TempDiskSpace))
		}
		if req.MinExp > 0 {
			sb.WriteString(fmt.Sprintf("&min_exp=%d", req.MinExp))
		}
		if req.MaxExp > 0 {
			sb.WriteString(fmt.Sprintf("&max_exp=%d", req.MaxExp))
		}
		
	case OpRegisterAssignment:
		req := pkt.(*RegisterAssignment)
		sb.WriteString("&t=ra")
		sb.WriteString("&g=" + url.QueryEscape(req.ComputerGUID))
		sb.WriteString(fmt.Sprintf("&c=%d&w=%d", req.CPUNum, req.WorkType))
		// Add work type specific parameters
		if req.WorkType == WorkTypeFactor {
			sb.WriteString(fmt.Sprintf("&n=%d&sf=%g", req.N, req.HowFarFactored))
			if req.FactorTo != 0.0 {
				sb.WriteString(fmt.Sprintf("&ef=%g", req.FactorTo))
			}
		} else if req.WorkType == WorkTypePFactor {
			sb.WriteString(fmt.Sprintf("&A=%.0f&b=%d&n=%d&C=%d&sf=%g&saved=%g",
				req.K, req.B, req.N, req.C, req.HowFarFactored, req.TestsSaved))
		} else if req.WorkType == WorkTypeFirstLL || req.WorkType == WorkTypeDblChk {
			sb.WriteString(fmt.Sprintf("&n=%d&sf=%g&p1=%d",
				req.N, req.HowFarFactored, req.HasBeenPMinus1ed))
		} else if req.WorkType == WorkTypePMinus1 {
			sb.WriteString(fmt.Sprintf("&A=%.0f&b=%d&n=%d&C=%d&B1=%d",
				req.K, req.B, req.N, req.C, req.B1))
			if req.B2 != 0 {
				sb.WriteString(fmt.Sprintf("&B2=%d", req.B2))
			}
		} else if req.WorkType == WorkTypePRP {
			sb.WriteString(fmt.Sprintf("&A=%.0f&b=%d&n=%d&C=%d&sf=%g",
				req.K, req.B, req.N, req.C, req.HowFarFactored))
			if req.PRPBase > 0 {
				sb.WriteString(fmt.Sprintf("&bases=%d", req.PRPBase))
			}
		}
		
	case OpAssignmentProgress:
		req := pkt.(*AssignmentProgress)
		sb.WriteString("&t=ap")
		sb.WriteString("&g=" + url.QueryEscape(req.ComputerGUID))
		sb.WriteString("&k=" + url.QueryEscape(req.AssignmentUID))
		if req.Stage != "" {
			sb.WriteString("&p=" + url.QueryEscape(req.Stage))
		}
		sb.WriteString(fmt.Sprintf("&e=%g", req.PctComplete))
		if req.CurrentIteration > 0 {
			sb.WriteString(fmt.Sprintf("&i=%d", req.CurrentIteration))
		}
		if req.Residue != "" {
			sb.WriteString("&r=" + url.QueryEscape(req.Residue))
		}
		if req.FFTLen > 0 {
			sb.WriteString(fmt.Sprintf("&fft=%d", req.FFTLen))
		}
		
	case OpAssignmentResult:
		req := pkt.(*AssignmentResult)
		sb.WriteString("&t=ar")
		sb.WriteString("&g=" + url.QueryEscape(req.ComputerGUID))
		sb.WriteString("&k=" + url.QueryEscape(req.AssignmentUID))
		sb.WriteString(fmt.Sprintf("&w=%d", req.WorkType))
		sb.WriteString(fmt.Sprintf("&A=%.0f&b=%d&n=%d&C=%d",
			req.K, req.B, req.N, req.C))
		if req.Residue != "" {
			sb.WriteString("&r=" + url.QueryEscape(req.Residue))
		}
		if req.ResidueType > 0 {
			sb.WriteString(fmt.Sprintf("&rt=%d", req.ResidueType))
		}
		if req.ErrorCode > 0 {
			sb.WriteString(fmt.Sprintf("&ec=%d", req.ErrorCode))
		}
		if req.ErrorMessage != "" {
			sb.WriteString("&em=" + url.QueryEscape(req.ErrorMessage))
		}
	}
	
	return sb.String()
}

// sendRequest sends HTTP GET request to PrimeNet server
func (c *Client) sendRequest(args string) (string, error) {
	requestURL := c.baseURL + "?" + args
	
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	var result strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			result.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
	
	return result.String(), nil
}

// parseResponse parses server response into packet structure
func (c *Client) parseResponse(resp string, operation int, result any) error {
	// Simple parsing - extract key-value pairs from response
	// Format: key=value\rkey=value\r...
	lines := strings.Split(resp, "\r")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Parse based on operation type
		switch operation {
		case OpGetAssignment:
			req := result.(*GetAssignment)
			switch key {
			case "k":
				req.AssignmentUID = value
			case "w":
				fmt.Sscanf(value, "%d", &req.WorkType)
			case "A":
				fmt.Sscanf(value, "%f", &req.K)
			case "b":
				fmt.Sscanf(value, "%d", &req.B)
			case "n":
				fmt.Sscanf(value, "%d", &req.N)
			case "C":
				fmt.Sscanf(value, "%d", &req.C)
			case "sf":
				fmt.Sscanf(value, "%f", &req.HowFarFactored)
			case "ef":
				fmt.Sscanf(value, "%f", &req.FactorTo)
			case "p1":
				fmt.Sscanf(value, "%d", &req.HasBeenPMinus1ed)
			case "B1":
				fmt.Sscanf(value, "%d", &req.B1)
			case "B2":
				fmt.Sscanf(value, "%d", &req.B2)
			case "CR":
				fmt.Sscanf(value, "%d", &req.Curves)
			case "bases":
				fmt.Sscanf(value, "%d", &req.PRPBase)
			case "rt":
				fmt.Sscanf(value, "%d", &req.PRPResidueType)
			case "dblchk":
				fmt.Sscanf(value, "%d", &req.PRPDblChk)
			case "squarings":
				fmt.Sscanf(value, "%d", &req.NumSquarings)
			case "known":
				req.KnownFactors = value
			}
		}
	}
	
	return nil
}

