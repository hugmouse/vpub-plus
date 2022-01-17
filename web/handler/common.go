// Code generated by go generate; DO NOT EDIT.

package handler

var TplCommonMap = map[string]string{
	"board_form": `{{ define "board_form" }}
<div class="field">
    <label for="forumId">Forum</label>
    <select name="forumId" id="forumId">
        {{ range .Forums }}
        <option value="{{ .Id }}" {{ if eq .Id $.ForumId }}selected{{ end }}>{{ .Name }}</option>
        {{ end }}
    </select>
</div>
<div class="field">
    <label for="name">Name</label>
    <input type="text" name="name" id="name" value="{{ .Name }}" autocomplete="off" maxlength="120" required autofocus/>
</div>
<div class="field">
    <label for="description">Description</label>
    <textarea class="editor" name="description" id="description" required>{{ .Description }}</textarea>
</div>
<div class="field">
    <label for="name">Position</label>
    <input type="number" name="position" id="position" value="{{ .Position }}" autocomplete="off" required/>
</div>
{{ end }}`,
	"forum_form": `{{ define "forum_form" }}
<div class="field">
  <label for="name">Name</label>
  <input type="text" name="name" id="name" value="{{ .Name }}" autocomplete="off" maxlength="120" required autofocus/>
</div>
<div class="field">
  <label for="name">Position</label>
  <input type="number" name="position" id="position" value="{{ .Position }}" autocomplete="off" required/>
</div>
{{ end }}`,
	"layout": `{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="/style.css"/>
    <title>{{ .settings.Name }}</title>
    {{ template "head" . }}
</head>
<body>
    <header>
        <span><a href="/">{{ .settings.Name }}</a></span>
        <nav>
            {{ if logged }}
            {{ if .hasNotifications }}<a href="/notifications" class="notifications">New replies</a> {{ end }} <a href="/account">{{ .logged.Name }}</a> (<a href="/logout">logout</a>)
            {{ else }}
            <a href="/login">login</a> <a href="/register">register</a>
            {{ end }}
        </nav>
    </header>
<!--    <p>{{ template "breadcrumb" . }}</p>-->
    {{ template "content" . }}
</body>
</html>
{{ end }}
{{ define "head" }}{{ end }}
{{ define "breadcrumb" }}{{ end }}`,
	"post_form": `{{ define "post_form" }}
<input type="hidden" name="topicId" value="{{ .TopicId }}">
<div class="field">
    <label for="subject">Subject</label>
    <input type="text" name="subject" id="subject" value="{{ .Subject }}" autocomplete="off" maxlength="115" required autofocus/>
</div>
<div class="field">
    <label for="content">Content</label>
    <textarea class="editor" name="content" id="content" required>{{ .Content }}</textarea>
</div>
{{ if .IsAdmin }}
<details>
    <summary>Admin options</summary>
    <div class="field">
        <label for="sticky">Sticky</label>
        <select name="sticky" id="sticky">
            <option value="false">false</option>
            <option value="true" {{ if .IsSticky }}selected{{ end }}>true</option>
        </select>
    </div>
    <div class="field">
        <label for="locked">Locked</label>
        <select name="locked" id="locked">
            <option value="false">false</option>
            <option value="true" {{ if .IsLocked }}selected{{ end }}>true</option>
        </select>
    </div>
</details>
{{ end }}
{{ end }}`,
	"posts": `{{ define "posts" }}
<!--{{ if . }}-->
<!--<ol class="posts">-->
<!--    {{ if . }}-->
<!--    {{ range . }}-->
<!--    <li>-->
<!--        <article>-->
<!--            <header><h2><a href="/posts/{{ .Id }}">{{ .Subject }}</a></h2> ({{ .Replies }})</header>-->
<!--            <div><a href="/~{{ .User }}">{{ .User }}</a>{{ if .Topic }} in <a href="/topics/{{ .Topic }}">{{ .Topic }}</a>{{ end }} {{ timeAgo .CreatedAt }}</div>-->
<!--        </article>-->
<!--    </li>-->
<!--    {{ end }}-->
<!--    {{ else }}-->
<!--    <li>No post yet</li>-->
<!--    {{ end }}-->
<!--</ol>-->
<!--{{ end }}-->

{{ if . }}
<table class="posts">
    <thead>
        <tr>
            <th>Subject</th>
            <th>Author</th>
            <th>Replies</th>
            <th>Updated</th>
        </tr>
    </thead>
    <tbody>
    {{ range . }}
    <tr>
        <td>
            <h2><a href="/posts/{{ .Id }}">{{ .Subject }}</a></h2>
        </td>
        <td style="text-align: center;"><a href="/~{{ .User }}">{{ .User }}</a></td>
        <td style="text-align: center;">{{ .Replies }}</td>
        <td style="text-align: center">{{ .DateUpdated }}</td>
    </tr>
    {{ end }}
    </tbody>
</table>
{{ else }}
<p>No post yet</p>
{{ end }}
{{ end }}`,
	"postsTopic": `{{ define "postsTopic" }}
{{ if . }}
<table class="posts">
    <thead>
    <tr>
        <th>Topic</th>
        <th>Subject</th>
        <th>Author</th>
        <th>Replies</th>
        <th>Updated</th>
    </tr>
    </thead>
    <tbody>
    {{ range . }}
    <tr>
        <td style="text-align: center;"><a href="/topics/{{ .Topic }}">{{ .Topic }}</a></td>
        <td>
            <h2><a href="/posts/{{ .Id }}">{{ .Subject }}</a></h2>
        </td>
        <td style="text-align: center;"><a href="/~{{ .User }}">{{ .User }}</a></td>
        <td style="text-align: center;">{{ .Replies }}</td>
        <td style="text-align: center">{{ .DateUpdated }}</td>
    </tr>
    {{ end }}
    </tbody>
</table>
{{ else }}
<p>No post yet</p>
{{ end }}
{{ end }}`,
}
