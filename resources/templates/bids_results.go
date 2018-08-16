package templates

var BidsResults = `
{{ define "content" }}
<br/><br/>
{{ .Badge }}

<h1>
    {{ .Header }}
</h1>
<br/>

<p>
    {{ range $val := .Issues.Errors }}
    <hr>
    <div>
        {{ $val.Severity }}, {{ $val.Key }}
    </div>
    <div>
        Reason: {{ $val.Reason }}
    </div>
    <div>
        {{ range $file := $val.Files }}
        <div>Filename: {{ $file.File.Name }} (Code: {{ $file.Code }})</div>
        <div>Path: {{ $file.File.Path }}</div>
        {{ end }}
    </div>
    {{ end }}

    {{ range $val := .Issues.Warnings }}
    <hr>
    <div>
        {{ $val.Severity }}, {{ $val.Key }}
    </div>
    <div>
        Reason: {{ $val.Reason }}
    </div>
    <div>
        {{ range $file := $val.Files }}
        <div>Filename: {{ $file.File.Name }} (Code: {{ $file.Code }})</div>
        <div>Path: {{ $file.File.Path }}</div>
        {{ end }}
    </div>
    {{ end }}
</p>

{{ if .Summary }}
<p>
    <div>Summary</div>
    <div>Sessions: {{ .Summary.Sessions }}</div>
    <div>Subjects: {{ .Summary.Subjects }}</div>
    <div>Tasks: {{ .Summary.Tasks }}</div>
    <div>Modalities: {{ .Summary.Modalities }}</div>
    <div>Total files: {{ .Summary.TotalFiles }}</div>
    <div>Size: {{ .Summary.Size }}</div>
</p>
{{ end }}

{{ end }}
`
