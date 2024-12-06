package main

import (
  // "errors"
  "fmt"
  "math"
  "os"
)

//--------------------------------------------------------------------
//Geometry
type Geometry interface {
  draw(scn screen) (err error) 
  shape() (s string)
}

type Color int

type Colors struct {
  Red   Color
  Green Color
  Blue  Color
}

type Point struct {
  x int
  y int
}

const (
  red    = 0
  green  = 1
  blue   = 2
  yellow = 3
  orange = 4
  purple = 5
  brown  = 6
  black  = 7
  white  = 8
)

type Triangle struct {
  pt0 Point
  pt1 Point
  pt2 Point
  c Color
}

type Rectangle struct {
  ll Point
  ur Point
  c Color
}

type Circle struct {
  cp Point
  r int
  c Color
}



//--------------------------------------------------------------------
//Screen


// display
// TODO: you must implement the struct for this variable, and the interface it implements (screen)



var display Display

//Declare screen interface
type screen interface {
  initialize(maxX, maxY int) 
  getMaxXY() (maxX, maxY int)
  drawPixel(x, y int, c Color) (err error) 
  getPixel(x, y int) (c Color, err error) 
  clearScreen() 
  screenShot(f string) (err error) 
}

type Display struct{
  maxX , maxY int
  matrix [][] Color
}

func (d *Display) initialize(maxX, maxY int){
  d.maxX = maxX
  d.maxY = maxY

  d.matrix = make([][] Color, maxY)
  for i := 0; i < maxY; i++ {
    d.matrix[i] = make([]Color, maxX)
    for j := 0; j < maxX; j++ {
      d.matrix[i][j] = white
    }
  }
}

func (d *Display) getMaxXY() (int, int){
  return d.maxX,d.maxY
}

func (d *Display) drawPixel (x, y int, c Color) (err error) {
  point := Point{x,y}  
  if colorUnknown(c){
    return colorUnknownErr{Message: "Color not valid!"}
  } else if outOfBounds(point, d){
    return outOfBoundsErr{Message: "Out of bounds error!"}
  } else {
    d.matrix[y][x] = c
  }
  return nil
}

func (d *Display) getPixel (x, y int) (c Color, err error) {
  point := Point{x,y}
  if outOfBounds(point, d){
    return -1, outOfBoundsErr{Message: "Out of bounds error!"}
  } else{
    return d.matrix[y][x], nil
  }
}

func (d *Display) clearScreen() {
  d.initialize(d.maxX, d.maxY)
}

func (d *Display) screenShot(f string) (err error) {
  cmap := []Colors{
    {255, 0, 0},
    {0, 255, 0},
    {0, 0, 255},
    {255, 255, 0},
    {255, 164, 0},
    {128, 0, 128},
    {165, 42, 42},
    {0, 0, 0},
    {255, 255, 255},
  }

  newFile := f + ".ppm"
  
  file, err := os.Create(newFile)
  if err != nil {
     os.Create(newFile)
  }

  _, err = fmt.Fprintf(file, "P3\n")
  if err != nil {
     return err
  }

  _, err = fmt.Fprintf(file, "%d %d\n", d.maxX, d.maxY)
  if err != nil {
     return err
  }

  _, err = fmt.Fprintf(file, "255\n")
  if err != nil {
     return err
  }

  for i:=0; i < d.maxY; i++ {
    for j:=0; j < d.maxX; j++ {
      _, err = fmt.Fprint(file, cmap[d.matrix[i][j]].Red)
      if err != nil {
         return err
      }

      _, err = fmt.Fprint(file," ")
      if err != nil {
         return err
      }
      _, err = fmt.Fprint(file, cmap[d.matrix[i][j]].Green)
      if err != nil {
         return err
      }

      _, err = fmt.Fprint(file," ")
      if err != nil {
         return err
      }
      _, err = fmt.Fprint(file, cmap[d.matrix[i][j]].Blue)
      if err != nil {
         return err
      }

      _, err = fmt.Fprint(file," ")
      if err != nil {
         return err
      }
    }
    _, err = fmt.Fprint(file, "\n")
    if err != nil {
       return err
    }
  }
  return nil
}



//--------------------------------------------------------------------
//Color Errors
type colorUnknownErr struct {
  Message string
}

func (colorError colorUnknownErr) Error() string {
  return colorError.Message
}

func colorUnknown(c Color) bool {
  if c < 0 || c > 8 {
    return true
  } else{
    return false
  }
}

//Bounds Errors
type outOfBoundsErr struct {
  Message string
}

func (boundsError outOfBoundsErr) Error() string {
  return boundsError.Message
}

func outOfBounds(p Point, s screen) bool {
  maxX, maxY := s.getMaxXY()

  if p.x > maxX || p.y > maxY || p.x < 0 || p.y < 0 {
    return true
  } else{
    return false
  }
}


//--------------------------------------------------------------------
//Geometry
// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate(l0, d0, l1, d1 int) (values []int) {
  a := float64(d1-d0) / float64(l1-l0)
  d := float64(d0)

  count := l1 - l0 + 1
  for ; count > 0; count-- {
    values = append(values, int(d))
    d = d + a
  }
  return
}

// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
  if outOfBounds(tri.pt0, scn) || outOfBounds(tri.pt1, scn) || outOfBounds(tri.pt2, scn) {
    return outOfBoundsErr{Message: "Out of bounds error!"}
  }
  if colorUnknown(tri.c) {
    return colorUnknownErr{Message: "Color not valid!"}
  }

  y0 := tri.pt0.y
  y1 := tri.pt1.y
  y2 := tri.pt2.y

  // Sort the points so that y0 <= y1 <= y2
  if y1 < y0 {
    tri.pt1, tri.pt0 = tri.pt0, tri.pt1
  }
  if y2 < y0 {
    tri.pt2, tri.pt0 = tri.pt0, tri.pt2
  }
  if y2 < y1 {
    tri.pt2, tri.pt1 = tri.pt1, tri.pt2
  }

  x0, y0, x1, y1, x2, y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

  x01 := interpolate(y0, x0, y1, x1)
  x12 := interpolate(y1, x1, y2, x2)
  x02 := interpolate(y0, x0, y2, x2)

  // Concatenate the short sides

  x012 := append(x01[:len(x01)-1], x12...)

  // Determine which is left and which is right
  var x_left, x_right []int
  m := len(x012) / 2
  if x02[m] < x012[m] {
    x_left = x02
    x_right = x012
  } else {
    x_left = x012
    x_right = x02
  }

  // Draw the horizontal segments
  for y := y0; y <= y2; y++ {
    for x := x_left[y-y0]; x <= x_right[y-y0]; x++ {
      scn.drawPixel(x, y, tri.c)
    }
  }
  return
}


//https://www.redblobgames.com/grids/circle-drawing/
func insideCircle(center, tile Point, radius int) bool {
  dx := center.x - tile.x
  dy := center.y - tile.y
  distance := math.Sqrt(float64(dx)*float64(dx) + float64(dy)*float64(dy))
  return distance <= float64(radius)
}

func (circ Circle) draw(scn screen) (err error) {
  leftX := circ.cp.x - circ.r
  leftPoint := Point{leftX, circ.cp.y}
  rightX := circ.cp.x + circ.r
  rightPoint := Point{rightX, circ.cp.y}
  upY := circ.cp.y + circ.r
  upPoint := Point{circ.cp.x, upY}
  downY := circ.cp.y - circ.r
  downPoint := Point{circ.cp.x, downY}
  
  if outOfBounds(leftPoint, scn) || outOfBounds(rightPoint, scn) || outOfBounds(downPoint, scn) || outOfBounds(upPoint, scn) {
    return outOfBoundsErr{Message: "Out of bounds error!"}
  }
  if colorUnknown(circ.c) {
    return colorUnknownErr{Message: "Color not valid!"}
  }
  // Draw the circle
  width, height := scn.getMaxXY()
  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      if insideCircle(circ.cp, Point{x, y}, circ.r) {
        scn.drawPixel(x, y, circ.c)
      }
    }
  }  
  return nil
}


func (rect Rectangle) draw(scn screen) (err error) {
  if outOfBounds(rect.ll, scn) || outOfBounds(rect.ur, scn) {
    return outOfBoundsErr{Message: "Out of bounds error!"}
  }
  if colorUnknown(rect.c) {
    return colorUnknownErr{Message: "Color not valid!"}
  }
  
  //Draw the rectangle
  // ul := Point{rect.ll.x, rect.ur.y}
  // lr := Point{rect.ur.x, rect. ll.y}
  for y := rect.ll.y; y < rect.ur.y; y++ {
    for x := rect.ll.x; x < rect.ur.x; x++ {
      scn.drawPixel(x, y, rect.c)
    }
  } 
  return nil
}


