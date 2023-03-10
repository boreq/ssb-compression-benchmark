# Results
```
goarch=amd64
goos=linux
cpu=AMD Ryzen 7 3700X 8-Core Processor
```
## Performance
### Few feed messages (compressing messages in batches of 10)
![](./Few&#32;feed&#32;messages&#32;\(compressing&#32;messages&#32;in&#32;batches&#32;of&#32;10\).png)
```
    Brotli (default) = 2.01 ratio
       Brotli (best) = 1.95 ratio
      Deflate (best) = 1.90 ratio
   Deflate (default) = 1.88 ratio
               LZMA2 = 1.85 ratio
   Deflate (fastest) = 1.85 ratio
         ZSTD (best) = 1.84 ratio
                LZMA = 1.82 ratio
      ZSTD (default) = 1.82 ratio
      ZSTD (fastest) = 1.80 ratio
    Brotli (fastest) = 1.63 ratio
              Snappy = 1.54 ratio
           S2 (best) = 1.53 ratio
         S2 (better) = 1.51 ratio
        S2 (default) = 1.50 ratio
```
### Few feed messages (compressing messages in batches of 100)
![](./Few&#32;feed&#32;messages&#32;\(compressing&#32;messages&#32;in&#32;batches&#32;of&#32;100\).png)
```
    Brotli (default) = 2.51 ratio
       Brotli (best) = 2.43 ratio
      Deflate (best) = 2.38 ratio
         ZSTD (best) = 2.37 ratio
      ZSTD (default) = 2.33 ratio
   Deflate (default) = 2.32 ratio
               LZMA2 = 2.32 ratio
                LZMA = 2.32 ratio
      ZSTD (fastest) = 2.30 ratio
   Deflate (fastest) = 2.26 ratio
    Brotli (fastest) = 2.09 ratio
           S2 (best) = 1.83 ratio
              Snappy = 1.82 ratio
         S2 (better) = 1.81 ratio
        S2 (default) = 1.77 ratio
```
### Few feed messages (compressing messages individually)
![](./Few&#32;feed&#32;messages&#32;\(compressing&#32;messages&#32;individually\).png)
```
       Brotli (best) = 1.28 ratio
    Brotli (default) = 1.25 ratio
      Deflate (best) = 1.21 ratio
               LZMA2 = 1.11 ratio
         ZSTD (best) = 1.10 ratio
   Deflate (default) = 1.10 ratio
      ZSTD (default) = 1.10 ratio
   Deflate (fastest) = 1.09 ratio
    Brotli (fastest) = 1.08 ratio
              Snappy = 1.04 ratio
      ZSTD (fastest) = 1.04 ratio
                LZMA = 0.99 ratio
           S2 (best) = 0.93 ratio
         S2 (better) = 0.93 ratio
        S2 (default) = 0.92 ratio
```
### Many feed messages (compressing messages in batches of 10)
![](./Many&#32;feed&#32;messages&#32;\(compressing&#32;messages&#32;in&#32;batches&#32;of&#32;10\).png)
```
    Brotli (default) = 1.91 ratio
       Brotli (best) = 1.84 ratio
      Deflate (best) = 1.79 ratio
   Deflate (default) = 1.77 ratio
         ZSTD (best) = 1.77 ratio
      ZSTD (default) = 1.75 ratio
               LZMA2 = 1.75 ratio
                LZMA = 1.73 ratio
      ZSTD (fastest) = 1.73 ratio
   Deflate (fastest) = 1.72 ratio
    Brotli (fastest) = 1.61 ratio
           S2 (best) = 1.43 ratio
        S2 (default) = 1.40 ratio
              Snappy = 1.39 ratio
         S2 (better) = 1.37 ratio
```
### Many feed messages (compressing messages in batches of 100)
![](./Many&#32;feed&#32;messages&#32;\(compressing&#32;messages&#32;in&#32;batches&#32;of&#32;100\).png)
```
    Brotli (default) = 2.48 ratio
       Brotli (best) = 2.41 ratio
         ZSTD (best) = 2.38 ratio
      Deflate (best) = 2.37 ratio
      ZSTD (default) = 2.34 ratio
      ZSTD (fastest) = 2.32 ratio
               LZMA2 = 2.31 ratio
                LZMA = 2.31 ratio
   Deflate (default) = 2.31 ratio
   Deflate (fastest) = 2.22 ratio
    Brotli (fastest) = 2.12 ratio
           S2 (best) = 1.82 ratio
         S2 (better) = 1.77 ratio
              Snappy = 1.77 ratio
        S2 (default) = 1.75 ratio
```
### Many feed messages (compressing messages individually)
![](./Many&#32;feed&#32;messages&#32;\(compressing&#32;messages&#32;individually\).png)
```
    Brotli (default) = 1.38 ratio
       Brotli (best) = 1.36 ratio
      Deflate (best) = 1.28 ratio
   Deflate (default) = 1.23 ratio
               LZMA2 = 1.20 ratio
         ZSTD (best) = 1.19 ratio
      ZSTD (default) = 1.19 ratio
   Deflate (fastest) = 1.12 ratio
      ZSTD (fastest) = 1.12 ratio
                LZMA = 1.11 ratio
    Brotli (fastest) = 1.07 ratio
              Snappy = 1.02 ratio
           S2 (best) = 1.00 ratio
        S2 (default) = 0.97 ratio
         S2 (better) = 0.96 ratio
```
