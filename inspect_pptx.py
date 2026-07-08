from pptx import Presentation
import sys

prs = Presentation(sys.argv[1])
for i, layout in enumerate(prs.slide_layouts):
    print(f"Layout {i}: {layout.name}")
    for j, ph in enumerate(layout.placeholders):
        print(f"  Placeholder {j}: {ph.name} (type: {ph.type})")
