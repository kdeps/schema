#!/bin/bash
# Fix pkldoc root index and create current symlink

set -e

PKLDOC_DIR="build/pkldoc/pkldoc"

if [ ! -d "$PKLDOC_DIR/schema.kdeps.com/core" ]; then
  echo "No pkldoc output found, skipping index fix"
  exit 0
fi

# Find the version directory (should be only one after a fresh generation)
VERSION_DIR=$(ls -1 "$PKLDOC_DIR/schema.kdeps.com/core" | grep -v "^current$" | head -1)

if [ -z "$VERSION_DIR" ]; then
  echo "No version directory found, skipping index fix"
  exit 0
fi

echo "Found version: $VERSION_DIR"

# Create 'current' symlink
echo "Creating 'current' symlink..."
cd "$PKLDOC_DIR/schema.kdeps.com/core"
ln -sf "$VERSION_DIR" current
cd - > /dev/null

# Update root index.html
echo "Updating root index.html..."
cat > "$PKLDOC_DIR/index.html" <<'EOF'
<!DOCTYPE html>
<html>
  <head>
    <title>Pkldoc</title>
    <script src="scripts/pkldoc.js" defer="defer"></script>
    <script src="scripts/scroll-into-view.min.js" defer="defer"></script>
    <link href="styles/pkldoc.css" media="screen" type="text/css" rel="stylesheet">
    <link rel="icon" type="image/svg+xml" href="images/favicon.svg">
    <link rel="apple-touch-icon" sizes="180x180" href="images/apple-touch-icon.png">
    <link rel="icon" type="image/png" sizes="32x32" href="images/favicon-32x32.png">
    <link rel="icon" type="image/png" sizes="16x16" href="images/favicon-16x16.png">
    <meta charset="UTF-8">
  </head>
  <body onload="onLoad()">
    <header>
      <div id="search"><i id="search-icon" class="material-icons">search</i><input id="search-input" type="search" placeholder="Click or press 'S' to search" autocomplete="off" data-root-url-prefix=""></div>
    </header>
    <main>
      <h1 id="declaration-title"></h1>
      <ul class="member-group-links">
        <li><a href="#_packages">Packages</a></li>
      </ul>
      <div class="member-group">
        <div id="_packages" class="anchor"> </div>
        <h2 class="member-group-title">Packages</h2>
        <ul>
          <li>
            <div id="schema.kdeps.com%2Fcore" class="anchor"> </div>
            <div class="member with-page-link"><a class="member-selflink material-icons" href="#schema.kdeps.com%2Fcore">link</a>
              <div class="member-left">
                <div class="member-modifiers">package </div>
              </div>
              <div class="member-main">
                <div class="member-signature"><a class="name-decl" href="./schema.kdeps.com/core/current/index.html">schema.kdeps.com/core</a></div>
                <div class="doc-comment"><p>Core PKL modules for KDEPS</p></div>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </main>
  </body>
</html>
EOF

# Update search index
echo "searchData='[{\"name\":\"schema.kdeps.com/core\",\"kind\":0,\"url\":\"schema.kdeps.com/core/current/index.html\"}]';" > "$PKLDOC_DIR/search-index.js"

echo "✅ Root index and symlink created for version $VERSION_DIR"
