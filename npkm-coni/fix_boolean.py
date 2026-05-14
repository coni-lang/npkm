with open('main.coni', 'r') as f:
    content = f.read()

target = """(if (boolean? changed-when-expr) changed-when-expr"""
replacement = """(if (or (= changed-when-expr true) (= changed-when-expr false)) changed-when-expr"""

content = content.replace(target, replacement)
with open('main.coni', 'w') as f:
    f.write(content)
