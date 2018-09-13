package templates

// TODO: Switch Login || Logout

// Layout is the main site template. It includes the header and footer and
// embeds the content for every other page.
var Layout = `
{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://web.gin.g-node.org/css/semantic-2.2.10.min.css">
    <link rel="stylesheet" href="https://web.gin.g-node.org/css/gogs.css">
    <title>G-Node GIN BIDS validation</title>
</head>
<body>
    <div class="following bar light">
        <div class="ui container">
            <div class="ui grid">
                <div class="column">
                    <div class="ui top secondary menu">
                        <a class="item brand" href="https://gin.g-node.org/">
                            <img class="ui mini image" src="https://gin.g-node.org/img/favicon.png">
                        </a>
                        <a class="item" href="https://gin.g-node.org/">Back to GIN</a>
                        <a class="item" href="/repos">Repositories</a>
                        <a class="item" href="/pubvalidate">One-time validation</a>
                        <a class="item" href="/login">Login</a>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div class="ui stackable middle very relaxed page grid">
        <div class="container main-container plainbox">
            {{ template "content" . }}
        </div>
    </div>
    <footer>
        <div class="following bar light">
            <div class="ui container">
                <div class="ui grid">
                    <div class="column">
                        <div class="ui top secondary menu center">
                            <a class="item brand" href="http://www.g-node.org">
                                <img class="ui mini image"
                                     src="https://projects.g-node.org/assets/gnode-bootstrap-theme/1.2.0-snapshot/img/gnode-icon-50x50-transparent.png"/>
                                © G-Node, 2018
                            </a>
                            <a class="item brand" href="https://web.gin.g-node.org/G-Node/Info/wiki/about" target="_blank">About</a>
                            <a class="item brand" href="https://web.gin.g-node.org/G-Node/Info/wiki/imprint" target="_blank">Imprint</a>
                            <a class="item brand" href="https://web.gin.g-node.org/G-Node/Info/wiki/contact" target="_blank">Contact</a>
                            <a class="item brand" href="https://web.gin.g-node.org/G-Node/Info/wiki/Terms+of+Use" target="_blank">Terms of Use</a>

                            <div class="ui supersmall">
                                Hosted by:
                                <img class="ui bmbf image" src="https://web.gin.g-node.org/img/lmu.png"/>
                            </div>
                            <div class="ui supersmall">
                                Funded by:
                                <img class="ui bmbf image" src="https://web.gin.g-node.org/img/bmbf.png"/>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </footer>
</body>
</html>
{{ end }}
`
