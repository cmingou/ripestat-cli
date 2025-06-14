# Documentation

This directory contains documentation assets and resources for the ripestat-cli project.

## Structure

```
docs/
├── README.md           # This file - documentation index
└── images/             # Visual assets for documentation
    └── architecture.svg # System architecture diagram
```

## Images

### architecture.svg
- **Description**: System architecture diagram showing the flow from CLI input through API calls to formatted output
- **Format**: SVG (Scalable Vector Graphics)
- **Usage**: Used in both English and Chinese README files
- **Generated**: Created from Mermaid diagram using mermaid-cli

## Adding New Documentation Assets

When adding new documentation assets:

1. **Images**: Place in `docs/images/`
2. **Diagrams**: Use Mermaid format when possible, convert to SVG for inclusion
3. **Screenshots**: Use PNG format, optimize for web
4. **Videos**: Place in `docs/videos/` (create directory as needed)

## Diagram Source

The architecture diagram was generated from this Mermaid source:

```mermaid
graph TD
    A[CLI Input] --> B[Argument Parser]
    B --> C{Input Type Detection}
    C -->|ASN| D[ASN Handler]
    C -->|IPv4| E[IPv4 Handler]
    C -->|IPv6| F[IPv6 Handler]
    C -->|Invalid| G[Error Handler]
    
    D --> H[RIPEstat API Client]
    E --> H
    F --> H
    
    H --> I[AS Overview API]
    H --> J[RIR API]
    H --> K[Prefix Routing API]
    H --> L[MaxMind GeoLite API]
    
    I --> M[Table Formatter]
    J --> M
    K --> M
    L --> M
    
    M --> N[Console Output]
    G --> N
    
    style A fill:#e1f5fe
    style N fill:#e8f5e8
    style H fill:#fff3e0
    style C fill:#fce4ec
```

To regenerate the diagram:
```bash
mmdc -i architecture.mmd -o docs/images/architecture.svg -t neutral -b white
```