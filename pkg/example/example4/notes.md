Cases

- heap vs stack
- hem pass by value hemde pass by pointer 'da da heap kaçmasına rağmen pass by pointer çok yavaş kalıyor? why?
- Peki random ürettiğim data sayısını 10_000 yaptığım zaman nede pass by reference ile pass by value arasındaki fark azalıyor?

userCount'u 10_000, 1000 ve 10 olarak verdiğimizde neden pass by value ve pass by pointer'ın hız farkı oluştuğu:

- Küçük userCount (ör: 10 veya 1000) değerlerinde, Go'nun slice'ları ve küçük veri bloklarını stack üzerinde tutabilme ve optimize edebilme şansı daha yüksektir. Pass by value ile slice'ın kendisi (header: pointer, len, cap) kopyalanır, veri ise çoğunlukla stack'te kalır ve hızlı erişilir.
- Pass by pointer ile ise, fonksiyona bir pointer geçirildiğinde, Go'nun escape analysis algoritması daha fazla veriyi heap'e kaçırabilir ve pointer dereferencing işlemleri ek bir maliyet oluşturur. Küçük veri miktarlarında bu ek maliyet, value'ya göre daha belirgin olur.
- Ayrıca, pointer ile çalışırken CPU cache locality avantajı kaybolabilir ve pointer dereferencing işlemleri CPU için ekstra yük oluşturur.

Büyük userCount (ör: 10_000) değerlerinde ise:

- Hem value hem de pointer ile slice'ın içindeki veri zaten heap'e kaçar (stack'te tutulamayacak kadar büyük). Bu durumda, value ile slice header'ı kopyalamak ile pointer'ı kopyalamak arasında neredeyse fark kalmaz. Çünkü asıl veri zaten heap'te ve erişim şekli aynıdır.
- Bu yüzden büyük veri miktarlarında pass by value ve pass by pointer arasındaki performans farkı azalır.

Özet: Küçük slice'larda pass by value daha hızlıdır çünkü stack optimizasyonları ve pointer dereferencing maliyeti yoktur. Büyük slice'larda ise ikisi de heap'e kaçar ve performans farkı kaybolur.
