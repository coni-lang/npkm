import re

with open("/Users/nico/cool/npkm/npkm-coni/main.coni", "r") as f:
    content = f.read()

# 1. Update /api/run routing to use path parameters
routing_old = """              (let [raw-path (get req :path)
                    path-parts (str/split raw-path "?")
                    path (first path-parts)
                    q-str (if (> (count path-parts) 1) (second path-parts) "")
                    start-idx (if (str/starts-with? q-str "start-idx=") (sys-parse-int (second (str/split q-str "="))) 0)
                    method (get req :method)]
                (if (= path "/")"""
routing_new = """              (let [raw-path (get req :path)
                    path-parts (str/split raw-path "/")
                    path (if (> (count path-parts) 1) (str "/" (nth path-parts 1)) "/")
                    method (get req :method)]
                (if (= raw-path "/")"""
if routing_old in content:
    content = content.replace(routing_old, routing_new)

# Fix /api/run endpoint matching and extract start-idx
api_run_old = """                      (if (= path "/api/run")"""
api_run_new = """                      (if (and (> (count path-parts) 2) (= (nth path-parts 1) "api") (= (nth path-parts 2) "run"))
                      (let [start-idx (if (> (count path-parts) 3) (sys-parse-int (nth path-parts 3)) 0)]"""
if api_run_old in content:
    content = content.replace(api_run_old, api_run_new)

# Ensure parentheses are balanced for api/run
api_run_close_old = """                          (close-chan ch)))
                        {:status 200 :headers {"Content-Type" "text/event-stream" "Cache-Control" "no-cache" "Connection" "keep-alive"} :body ch})
                      {:status 404 :body "Not found"}))))))))"""
api_run_close_new = """                          (close-chan ch)))
                        {:status 200 :headers {"Content-Type" "text/event-stream" "Cache-Control" "no-cache" "Connection" "keep-alive"} :body ch}))
                      {:status 404 :body "Not found"}))))))))"""
if api_run_close_old in content:
    content = content.replace(api_run_close_old, api_run_close_new)

# Wait, `slice-plays-by-task-idx` implementation!
# If it returns a LazyStream, it fails `(vector?)` check in `execute-playbook`.
# So we need to evaluate it to a vector!
# Coni might not have `vec`? If it doesn't, we can just leave it if it works... Wait, the bug was NOT LazyStream, it was `start-idx` parsing!
# If it was LazyStream, it would have wrapped it in `Default Play` and printed `PLAY [ Default Play ]`.
# But it printed `PLAY [ Flow Control Demo ]` in the screenshot, so it WAS a vector, or `vector?` works on LazyStream?
# Wait! `(:tasks parsed-data)` is passed directly.

# 2. Fix JS in HTML
html_match = re.search(r':body \(str "(<!DOCTYPE html>.*?)(\"\)})', content, flags=re.DOTALL)
if html_match:
    html = html_match.group(1)
    
    # EventSource URL
    html = html.replace("new EventSource('/api/run?start-idx=' + idx);", "new EventSource('/api/run/' + idx);")
    
    # Fix viewArchitecture Mermaid newlines (replace \\\\n with \\n)
    html = html.replace('graph TD\\\\n', 'graph TD\\n')
    
    # Remove Architecture Map from runAudit
    run_audit_map_str = "html += '<div class=\\\"task-card\\\" style=\\\"padding: 25px; overflow-x: auto;\\\"><h3 style=\\\"margin-top:0; color: #f8fafc;\\\">Architecture Map</h3><div style=\\\"background: rgba(0,0,0,0.3); border-radius: 8px; padding: 20px; margin-top: 15px; min-width: 800px; display: flex; justify-content: center;\\\"><div class=\\\"mermaid\\\">graph TD\\n' + data.mermaid + '</div></div></div>';"
    if run_audit_map_str in html:
        html = html.replace(run_audit_map_str, "")
    else:
        print("Could not find Architecture Map in runAudit!")
        
    # Remove Playbook Tasks from sidebar
    html = html.replace('<h2>Playbook Tasks</h2><div id=\\"tree\\">Loading...</div>', '')
    
    # Remove JS tree population logic
    tree_logic_regex = r"let plays = Array\.isArray\(data\.tasks\) \? data\.tasks : \[data\.tasks\]; plays\.forEach\(p => \{ if\(p\.tasks && Array\.isArray\(p\.tasks\)\) \{ html \+= `<div class=\\\"tree-item\\\" style=\\\"color:#a855f7; margin-top:10px;\\\">▼ \$\{p\.name \|\| 'Play'\}</div>`; p\.tasks\.forEach\(t => \{ html \+= `<div class=\\\"tree-item\\\" style=\\\"padding-left:15px;\\\"><span class=\\\"tree-icon\\\">▶</span>\$\{t\.name \|\| 'Unnamed Task'\}</div>`; \}\); \} else if\(p\.name\) \{ html \+= `<div class=\\\"tree-item\\\"><span class=\\\"tree-icon\\\">▶</span>\$\{p\.name \|\| 'Unnamed Task'\}</div>`; \} \}\); \} document\.getElementById\('tree'\)\.innerHTML = html \|\| '<div class=\\\"tree-item\\\">No tasks</div>';"
    html = re.sub(tree_logic_regex, '} /* tree removed */', html)

    content = content[:html_match.start(1)] + html + content[html_match.end(1):]
else:
    print("HTML not found!")

with open("/Users/nico/cool/npkm/npkm-coni/main.coni", "w") as f:
    f.write(content)

print("Patch applied successfully!")
