import sys, re
def get_imbalance(s):
    stack = []
    pairs = {')': '(', ']': '[', '}': '{'}
    line_num = 1
    for idx, c in enumerate(s):
        if c == '\n': line_num+=1
        elif c in '([{': stack.append((c, line_num))
        elif c in ')]}':
            if not stack: return f"Extra {c} at line {line_num}"
            top, _ = stack.pop()
            if top != pairs[c]: return f"Mismatch {c} for {top} at line {line_num}"
    return f"Unclosed count: {len(stack)}" if stack else "OK"

with open('main.coni') as f: text = f.read()

def remove_strings(s):
    res = []
    i = 0
    while i < len(s):
        if s[i] == '"':
            res.append(' ')
            i += 1
            while i < len(s):
                if s[i] == '"' and s[i-1] != '\\':
                    res.append(' ')
                    i += 1
                    break
                elif s[i] == '\n':
                    res.append('\n')
                else:
                    res.append(' ')
                i += 1
        else:
            res.append(s[i])
            i += 1
    return "".join(res)

s_no_strings = remove_strings(text)
s_no_comments = re.sub(r';.*', '', s_no_strings)
print(get_imbalance(s_no_comments))
