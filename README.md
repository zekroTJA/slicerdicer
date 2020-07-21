# slicerdicer

A little command line tool to slice an image into same-sized parts. This was concipated to slice images into multipart Discord emotes.

## Example

Imput image: 

![](testdata/source/pog.jpg)

```
$ go run cmd/main.go \
    -i testdata/source/pog.jpg \
    -o testdata/results \
    -s 2 \
    -oname pog
```

Results:

![](testdata/results/pog_0_0.png) ![](testdata/results/pog_1_0.png)
<br/>
![](testdata/results/pog_0_1.png) ![](testdata/results/pog_1_1.png)

---

© 2020 Ringo Hoffmann (zekro Development)  
Covered by the MIT Licence.