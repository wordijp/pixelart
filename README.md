# pixelart

Convert calendar graph SVG to pixel art

source Graph SVG for Calendar(for example, GitHub and [Pixela](https://pixe.la))

<img src="./doc/img/vim-pixela.png" width="180" alt="vim" title="vim">

and Image, Convert to pixel art

- <img src="./doc/img/vim.png" width="80" alt="vim" title="vim">
 →
<img src="./doc/img/dst-vim-pixela.png" width="140" alt="vim" title="vim">
- <img src="./doc/img/grass.png" width="180" alt="vim" title="vim">
 →
<img src="./doc/img/dst-grass.png" width="200" alt="vim" title="vim">

this PSD image

- <img src="./doc/img/calendar.png" width="300" alt="vim" title="vim">
 →
<img src="./doc/img/dst-calendar.png" width="360" alt="vim" title="vim">


## Installation

```
$ go get github.com/wordijp/pixelart
```

## Requirement

- [ajstarks/svgo](https://github.com/ajstarks/svgo)
- [vmihailenco/msgpack](https://github.com/vmihailenco/msgpack)
- [wordijp/svgparser](https://github.com/wordijp/svgparser)
	- forked) [JoshVarga/svgparser](https://github.com/JoshVarga/svgparser)
- [oov/psd](https://github.com/oov/psd)

## Usage


```go
import "github.com/wordijp/pixelart/graph"
import "github.com/wordijp/pixelart/dot"

func ExampleConvertPrint() {
	var g graph.Data
	{
		file, _ := os.Open("calendar-graph.svg")
		defer file.Close()
		g, _ = graph.ParseCalendarGraphSvg(file)
	}
	
	var d dot.Data
	{
		file, _ := os.Open("dot-vim.png")
		defer file.Close()
		d, _ = dot.ParseDotPng(file)
	}
	
	buf := bytes.NewBuffer(nil)
	d.Convert(g).WriteSvgString(buf)

	fmt.Println(buf.String())
}
```

and see example

## License

The MIT License.
