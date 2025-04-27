# Go Bottleneck Demonstrasyonu

Bu proje, Go'daki performans darboğazlarını gösteren bir POC (Proof of Concept) çalışmasıdır.

## Proje Hakkında

Bu projede gösterilen temel noktalar:

- Ardışık (Sequential) ve eşzamanlı (Concurrent) kod arasındaki performans farkı
- Goroutine'lerde bellek sızıntısı (leak) oluşması
- Gerçek hayatta karşılaşılabilecek örnekler
- Micro ölçekte benchmark testleri ile performans izleme
- Macro ölçekte performans testi

## Kurulum ve Çalıştırma

```bash
# Projeyi klonlayın
git clone https://github.com/idagdelen/go-bottlenecks.git
cd go-bottlenecks

# Bağımlılıkları yükleyin
go mod download

# Uygulamayı çalıştırın
go run cmd/main.go
```

Uygulama http://localhost:8080 adresinde çalışacaktır.


## Swagger Dokümantasyonu

API dokümantasyonuna Swagger UI üzerinden erişebilirsiniz:

```
http://localhost:8080/swagger/index.html
```

### Swagger Dokümantasyonunu Güncelleme

Kodunuzda API endpoint'lerini veya parametrelerini değiştirdikten sonra Swagger dokümantasyonunu güncellemek için aşağıdaki adımları izleyin:

1. Öncelikle, swag CLI aracının yüklü olduğundan emin olun:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Swagger dökümanlarını yeniden oluşturun:
   ```bash
   swag init -g cmd/main.go -o pkg/docs
   ```

3. Uygulamayı yeniden başlatın:
   ```bash
   go run cmd/main.go
   ```

## Lisans

Bu proje [MIT Lisansı](LICENSE) altında lisanslanmıştır. 
