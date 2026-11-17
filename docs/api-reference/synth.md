---
description: SynthRequest, SynthResult, and App.Synthesize
icon: waveform-lines
---

# Synth

```go
type SynthRequest struct {
    VoiceRef    string  // character:voice (mutex with VoiceFile)
    VoiceFile   []byte  // raw Character YAML bytes
    VoiceName   string  // required with VoiceFile
    VoiceBase   string  // fallback character name when YAML has none
    Text        string
    Version     int
    SpeedFactor float64
    Output      string
    Preview     bool
    ReuseCache  *bool   // nil → use config synth.reuse_cache
}

type SynthResult struct {
    WAV       []byte
    Path      string
    Wrote     bool
    FromCache bool
}

func (a *App) Synthesize(ctx context.Context, req SynthRequest) (SynthResult, error)
```

{% tabs %}
{% tab title="By voice ref" %}
```go
res, err := app.Synthesize(ctx, api.SynthRequest{
    VoiceRef: "remilia:default",
    Text:     "こんにちは",
})
```
{% endtab %}

{% tab title="By YAML bytes" %}
```go
res, err := app.Synthesize(ctx, api.SynthRequest{
    VoiceFile: yamlBytes,
    VoiceName: "default",
    Text:      "こんにちは",
    Preview:   true,
})
```
{% endtab %}

{% tab title="Cache override" %}
```go
reuse := false
res, err := app.Synthesize(ctx, api.SynthRequest{
    VoiceRef:   "remilia:default",
    Text:       "こんにちは",
    ReuseCache: &reuse,
})
```
{% endtab %}
{% endtabs %}

{% hint style="info" %}
Write policy, filename format, and cache hashing match the CLI — see [Configuration](https://app.gitbook.com/s/A67wXyFHjvGTd6AJa6xT/reference/configuration) and [synth](https://app.gitbook.com/s/A67wXyFHjvGTd6AJa6xT/cli/synth).
{% endhint %}
