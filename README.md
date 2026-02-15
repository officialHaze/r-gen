# R-Gen — Report Generation Engine

## Overview

**R-Gen** is a lightweight **PDF report generation engine** designed to produce deterministic, print-ready PDFs directly from **HTML templates**. It works by injecting dynamic data into templates filled with placeholders, converting the rendered HTML into PDFs, and optionally **merging multiple PDFs into a single consolidated document**.

R-Gen also supports **standalone PDF generation**, where an HTML template can be uploaded and converted to a PDF in real time. This eliminates the overhead of first uploading a template, then injecting data, and finally generating the PDF, making it suitable for quick, one-off report generation scenarios.

The engine provides built-in capabilities to **preview** and **download** generated reports, making it suitable for both automated systems and manual review workflows. It operates without any database dependency or user registration requirements. All templates and generated PDFs are stored locally and segregated using **IP-based directory structures**.

R-Gen runs as a server exposed through a **REST API**, fully packaged and deployed using **Docker and Docker Compose**, ensuring consistent behavior across environments without dependency conflicts.

---

## Primary Motive

R-Gen was originally developed to support **SECVEIL’s report generation requirements**. During evaluation, existing solutions were found to be either overly complex, larger than the actual use case, or commercially restrictive.

The goal was to build a **minimalistic and efficient alternative** that fulfills the exact needs of the platform without unnecessary overhead. While its initial purpose was SECVEIL integration, R-Gen has been kept generic so that others can also try and use it. Additional features and enhancements may be introduced later if required.

---

## Target Use Cases

R-Gen is suitable for generating structured PDF reports such as:

- Security reports
- Audit reports
- Compliance reports
- Internal dashboards

---

## Core Strengths

- Template-based generation using HTML
- Standalone HTML-to-PDF generation without prior template storage
- Deterministic output for consistent report rendering
- HTML to PDF conversion pipeline
- Config-driven layouts
- Performance-focused design
- PDF merging capability
- Preview and download support
- No database dependency
- No registration or authentication requirements
- IP-based storage segregation
- Fully containerized deployment

---

## Key Focus Areas

R-Gen prioritizes:

- Consistent layout across generated reports
- Print-safe PDF output
- Support for custom branding via templates
- Automation readiness through REST APIs
- Developer-first usage
- Dynamic data injection into templates
- Easy integration with external systems
- Open-source accessibility

---

## Templating-First Design Philosophy

R-Gen is **heavily templating-driven** by design. This means the **data structure and the HTML template must be intentionally designed to work together**. To fully leverage the engine’s capabilities, users are expected to prepare and adhere to a consistent data format that aligns with their templates.

There is **no enforced schema** by R-Gen. The responsibility of defining and maintaining compatibility between data and templates lies entirely with the user.

---

## Example: Data Format (Dummy)

```json
{
  "device": {
    "name": "ABC-Firewall-01",
    "vendor": "ExampleVendor",
    "model": "FW-X1000"
  },
  "usage": {
    "cpu_percent": 35,
    "ram_percent": 42
  },
  "compliance": {
    "overall_percent": 82,
    "status": "Compliant"
  }
}
```

## Example: HTML Format (as per designed Data Format)

```html
<!doctype html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Firewall Report</title>
  </head>
  <body>
    <h1>Firewall Security Report</h1>

    <p><strong>Device Name:</strong> {{.device.name}}</p>
    <p><strong>Vendor:</strong> {{.device.vendor}}</p>
    <p><strong>Model:</strong> {{.device.model}}</p>

    <h2>Resource Usage</h2>
    <p>CPU Usage: {{.usage.cpu_percent}}%</p>
    <p>RAM Usage: {{.usage.ram_percent}}%</p>

    <h2>Compliance Summary</h2>
    <p>Status: {{.compliance.status}}</p>
    <p>Overall Compliance: {{.compliance.overall_percent}}%</p>
  </body>
</html>
```

---

## Important Notes

- Template placeholders must directly match the structure of the provided data.
- Nested objects, conditional rendering, loops, and formatting are handled at the template level.
- Poorly aligned data and templates may result in incomplete or incorrect reports.
- R-Gen intentionally avoids abstracting this layer to retain flexibility and performance.
