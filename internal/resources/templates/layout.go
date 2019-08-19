package templates

// TODO: Switch Login || Logout

// Layout is the main site template. It includes the header and footer and
// embeds the content for every other page.
var Layout = `
{{ define "layout" }}
<html>
	<!DOCTYPE html>
	<head data-suburl="">
		<link rel="shortcut icon" href="https://gin.g-node.org/img/favicon.png" />
		<link rel="stylesheet" href="https://gin.g-node.org/assets/font-awesome-4.6.3/css/font-awesome.min.css">
		<link rel="stylesheet" href="https://gin.g-node.org/assets/octicons-4.3.0/octicons.min.css">
		<link rel="stylesheet" href="https://gin.g-node.org/css/semantic-2.3.1.min.css">
		<link rel="stylesheet" href="https://gin.g-node.org/css/gogs.css?v=921e73e55b4d707a9a72151df987dce1">
		<title>Sign In - GIN Valid</title>
		<link rel="stylesheet" href="https://gin.g-node.org/css/custom.css">
		<meta name="twitter:card" content="summary" />
		<meta name="twitter:site" content="@gnode" />
		<meta name="twitter:title" content="GIN Valid"/>
		<meta name="twitter:description" content="G-Node GIN Validation service"/>
		<meta name="twitter:image" content="https://gin.g-node.org/img/favicon.png" />
	</head>
	<body>
		<div class="full height">
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
			{{ template "content" . }}
		</div>
		<footer>
			<div class="ui container">
				<div class="ui center links item brand footertext">
					<a href="http://www.g-node.org"><img class="ui mini footericon" src="https://projects.g-node.org/assets/gnode-bootstrap-theme/1.2.0-snapshot/img/gnode-icon-50x50-transparent.png"/>© G-Node, 2016-2019</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/about">About</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/imprint">Imprint</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/contact">Contact</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/Terms+of+Use">Terms of Use</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/Datenschutz">Datenschutz</a>
				</div>
				<div class="ui center links item brand footertext">
					<span>Powered by:      <a href="https://github.com/gogits/gogs"><img class="ui mini footericon" src="https://gin.g-node.org/img/gogs.svg"/></a>         </span>
					<span>Hosted by:       <a href="http://neuro.bio.lmu.de"><img class="ui mini footericon" src="https://gin.g-node.org/img/lmu.png"/></a>          </span>
					<span>Funded by:       <a href="http://www.bmbf.de"><img class="ui mini footericon" src="https://gin.g-node.org/img/bmbf.png"/></a>         </span>
					<span>Registered with: <a href="http://doi.org/10.17616/R3SX9N"><img class="ui mini footericon" src="https://gin.g-node.org/img/re3.png"/></a>          </span>
					<span>Recommended by:  <a href="https://www.nature.com/sdata/policies/repositories#neurosci"><img class="ui mini footericon" src="https://gin.g-node.org/img/sdatarecbadge.jpg"/><a href="https://journals.plos.org/plosone/s/data-availability#loc-neuroscience"><img class="ui mini footericon" src="https://gin.g-node.org/img/sm_plos-logo-sm.png"/></a></span>
				</div>
			</div>
		</footer>
	</body>
</html>
{{ end }}
`
