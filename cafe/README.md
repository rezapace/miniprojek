sebelum memulai jalankan
go mod init
go mod tidy

# User

- [x] user dapat registrasi dengan role sebagai customer atau kasir
- [x] user dapat login dengan email dan password
- [x] user dapat melihat daftar makanan
- [x] user dapat menambahkan pesanan baru (order)

# Kasir

- [x] kasir dapat menambahkan produk makanan baru
- [x] kasir dapat melihat detail pesanan berdasarkan ID pesanan
- [x] kasir dapat mengupdate order dalam sistem
- [x] kasir dapat menghapus order dalam sistem
- [x] kasir dapat mengubah data makanan berdasarkan ID makanan
- [x] kasir dapat menghapus data makanan berdasarkan ID makanan
- [x] kasir dapat melihat status pesanan berdasarkan ID pesanan

-t : tag atau memberikan nama pada image yang akan dibuat
nama_image : nama yang diberikan pada image
tag : versi dari image yang dibuat
. : menunjukkan bahwa Dockerfile yang digunakan berada pada direktori yang sama dengan terminal saat ini
docker build -t cafe-app:v1.0 .

FIXME:

- ketika mengunakan jwt output nya selalu error di sarankan mengganti jwt dan echo versi
- ketika di build docker tidak bisa include gorm terjadi erro
- pada routes belum di tambahkan user role nya berdasarkan siapa saja yang dapat mengakses di karenakan echo dan jwt nya
