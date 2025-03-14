package disease

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// Helper function to generate the PDF
func GeneratePrescriptionPDF(p *Prescription, output *bytes.Buffer) error {
	// PDF constants
	const (
		pageWidth     = 210.0
		margin        = 15.0
		contentWidth  = pageWidth - 2*margin
		leftColWidth  = 30.0
		rightColStart = pageWidth/2 + 10
		signatureY    = 250.0
	)

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()

	// Header section
	if p.HospitalLogo != "" {
		pdf.Image(p.HospitalLogo, margin, margin, 25, 0, false, "", 0, "")
	} else {
		pdf.SetFillColor(220, 220, 220)
		pdf.Rect(margin, margin, 25, 25, "F")
	}

	// Hospital info
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 102, 204)
	pdf.Cell(40, 10, p.HospitalName)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.Text(margin+30, margin+15, p.HospitalAddress)
	pdf.Text(margin+30, margin+22, p.HospitalPhone)

	pdf.SetDrawColor(0, 102, 204)
	pdf.SetLineWidth(0.5)
	pdf.Line(margin, margin+30, pageWidth-margin, margin+30)

	// Prescription title
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(0, 102, 204)
	pdf.Text(pageWidth/2-25, margin+45, "PRESCRIPTION")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.Text(margin, margin+45, fmt.Sprintf("No: %s", p.ID))
	pdf.Text(pageWidth-margin-50, margin+45, fmt.Sprintf("Date: %s", p.PrescribedDate.Format("02/01/2006")))

	// Patient information
	yPos := margin + 60
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 102, 204)
	pdf.Text(margin, yPos, "PATIENT INFORMATION")

	yPos += 10
	pdf.SetFillColor(240, 240, 255)
	pdf.Rect(margin, yPos, contentWidth, 35, "F")

	// Patient details
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(60, 60, 60)
	pdf.Text(margin+5, yPos+10, "Name:")
	pdf.Text(margin+5, yPos+20, "Gender:")
	pdf.Text(margin+5, yPos+30, "Address:")

	pdf.SetFont("Arial", "", 10)
	pdf.Text(margin+leftColWidth, yPos+10, p.PatientName)
	pdf.Text(margin+leftColWidth, yPos+20, p.PatientGender)

	pdf.SetFont("Arial", "B", 10)
	pdf.Text(rightColStart, yPos+10, "Age:")
	pdf.Text(rightColStart, yPos+20, "Patient ID:")

	pdf.SetFont("Arial", "", 10)
	pdf.Text(rightColStart+30, yPos+10, fmt.Sprintf("%d", p.PatientAge))

	// Diagnosis section
	yPos += 45
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 102, 204)
	pdf.Text(margin, yPos, "DIAGNOSIS")

	yPos += 8
	pdf.SetFillColor(245, 245, 245)
	pdf.Rect(margin, yPos, contentWidth, 15, "F")

	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(60, 60, 60)
	pdf.Text(margin+5, yPos+10, p.Diagnosis)

	// Medicines section
	yPos += 25
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 102, 204)
	pdf.Text(margin, yPos, "TREATMENT MEDICINES")

	yPos += 10
	colWidths := []float64{60, 30, 35, 30, 45}
	colNames := []string{"Medicine Name", "Dosage", "Frequency", "Duration", "Instructions"}

	// Table header
	pdf.SetFillColor(0, 102, 204)
	pdf.SetTextColor(255, 255, 255)

	xPos := margin
	for i, width := range colWidths {
		pdf.Rect(xPos, yPos, width, 10, "F")
		pdf.Text(xPos+3, yPos+7, colNames[i])
		xPos += width
	}

	// Table data
	yPos += 10
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(60, 60, 60)

	for i, medicine := range p.Medicines {
		if i%2 == 0 {
			pdf.SetFillColor(240, 240, 255)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.Rect(margin, yPos, contentWidth, 10, "F")

		xPos = margin
		pdf.Text(xPos+3, yPos+7, medicine.MedicineName)
		xPos += colWidths[0]
		pdf.Text(xPos+3, yPos+7, medicine.Dosage)
		xPos += colWidths[1]
		pdf.Text(xPos+3, yPos+7, medicine.Frequency)
		xPos += colWidths[2]
		pdf.Text(xPos+3, yPos+7, medicine.Duration)
		xPos += colWidths[3]
		pdf.Text(xPos+3, yPos+7, medicine.Notes)

		yPos += 10
	}

	// Notes section (if any)
	if p.Notes != "" {
		yPos += 10
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(0, 102, 204)
		pdf.Text(margin, yPos, "NOTES")

		yPos += 8
		pdf.SetFillColor(255, 250, 240)
		pdf.Rect(margin, yPos, contentWidth, 25, "F")

		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(60, 60, 60)
		pdf.SetXY(margin+5, yPos+5)
		pdf.MultiCell(contentWidth-10, 5, p.Notes, "", "", false)
	}

	// Doctor signature
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(80, 80, 80)

	dateText := fmt.Sprintf("Date %d/%d/%d",
		p.PrescribedDate.Day(),
		p.PrescribedDate.Month(),
		p.PrescribedDate.Year())

	pdf.Text(pageWidth-margin-70, signatureY, dateText)
	pdf.SetFont("Arial", "B", 10)
	pdf.Text(pageWidth-margin-60, signatureY+10, "Prescribing Doctor")

	pdf.SetDrawColor(0, 102, 204)
	pdf.SetLineWidth(0.2)
	pdf.Line(pageWidth-margin-80, signatureY+40, pageWidth-margin-20, signatureY+40)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.Text(pageWidth-margin-70, signatureY+45, p.DoctorName)

	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(80, 80, 80)
	pdf.Text(pageWidth-margin-70, signatureY+52, p.DoctorTitle)

	// Footer
	pdf.SetDrawColor(0, 102, 204)
	pdf.SetLineWidth(1)
	pdf.Line(margin, 280, pageWidth-margin, 280)

	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(80, 80, 80)
	pdf.Text(margin, 287, fmt.Sprintf("Contact: %s | %s", p.HospitalPhone, p.HospitalAddress))

	// Output PDF to buffer
	return pdf.Output(output)
}
