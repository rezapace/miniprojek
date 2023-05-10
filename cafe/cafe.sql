-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Waktu pembuatan: 10 Bulan Mei 2023 pada 01.18
-- Versi server: 10.4.25-MariaDB
-- Versi PHP: 8.1.10

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `cafe`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `food`
--

CREATE TABLE `food` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `price` decimal(10,2) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `food`
--

INSERT INTO `food` (`id`, `name`, `description`, `price`) VALUES
(2, 'Nasi Goreng', 'Nasi goreng istimewa dengan telur, ayam, dan bumbu rempah pilihan.', '15000.00');

-- --------------------------------------------------------

--
-- Struktur dari tabel `foods`
--

CREATE TABLE `foods` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext DEFAULT NULL,
  `description` longtext DEFAULT NULL,
  `price` double DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `foods`
--

INSERT INTO `foods` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `description`, `price`) VALUES
(1, '2023-05-09 21:36:25.756', '2023-05-09 21:39:17.907', '2023-05-09 21:40:24.845', 'Nasi Goreng Spesial', 'Nasi goreng dengan telur, ayam, sayuran, dan rempah-rempah khas Indonesia.', 225000),
(2, '2023-05-09 21:41:44.363', '2023-05-10 02:37:16.811', NULL, 'Nasi Goreng Spesial', 'Nasi goreng dengan telur, ayam, sayuran, dan rempah-rempah khas Indonesia.', 225000),
(3, '2023-05-09 21:41:51.319', '2023-05-09 21:41:51.319', NULL, 'Nasi Goreng 1', 'Nasi goreng istimewa dengan telur, ayam, dan bumbu rempah pilihan.', 151000),
(4, '2023-05-09 21:41:59.148', '2023-05-09 21:41:59.148', NULL, 'Nasi Goreng 2', 'Nasi goreng istimewa dengan telur, ayam, dan bumbu rempah pilihan.', 1521000),
(5, '2023-05-10 02:11:39.996', '2023-05-10 02:11:39.996', NULL, 'Nasi Goreng', 'Nasi goreng istimewa dengan telur, ayam, dan bumbu rempah pilihan.', 15000),
(6, '2023-05-10 02:11:48.108', '2023-05-10 02:11:48.108', NULL, 'Nasi Goreng', 'Nasi goreng istimewa dengan telur, ayam, dan bumbu rempah pilihan.', 15000),
(7, '2023-05-10 02:12:44.654', '2023-05-10 02:12:44.654', NULL, 'Nasi Goreng', 'Nasi goreng istimewa dengan telur, ayam, dan bumbu rempah pilihan.', 15000),
(8, '2023-05-10 02:37:14.482', '2023-05-10 02:37:14.482', NULL, 'Nasi Goreng', 'Nasi goreng istimewa dengan telur, ayam, dan bumbu rempah pilihan.', 15000);

-- --------------------------------------------------------

--
-- Struktur dari tabel `orders`
--

CREATE TABLE `orders` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `total_price` decimal(10,2) NOT NULL,
  `status` enum('proses','selesai') NOT NULL,
  `order_time` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `orders`
--

INSERT INTO `orders` (`id`, `user_id`, `total_price`, `status`, `order_time`) VALUES
(2, 1, '2000.00', 'proses', '2023-05-09 17:22:52'),
(3, 2, '2000.00', 'proses', '0000-00-00 00:00:00'),
(4, 2, '2000.00', 'proses', '0000-00-00 00:00:00'),
(5, 2, '2000.00', 'proses', '0000-00-00 00:00:00'),
(6, 2, '0.00', 'proses', '0000-00-00 00:00:00');

-- --------------------------------------------------------

--
-- Struktur dari tabel `order_details`
--

CREATE TABLE `order_details` (
  `id` int(11) NOT NULL,
  `order_id` int(11) NOT NULL,
  `food_id` int(11) NOT NULL,
  `quantity` int(11) NOT NULL,
  `price` decimal(10,2) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `order_details`
--

INSERT INTO `order_details` (`id`, `order_id`, `food_id`, `quantity`, `price`) VALUES
(2, 2, 2, 2, '1000.00'),
(3, 3, 2, 2, '1000.00'),
(4, 4, 2, 2, '1000.00'),
(5, 5, 2, 2, '1000.00');

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `name` longtext DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `password` longtext DEFAULT NULL,
  `userrole` longtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `name`, `email`, `password`, `userrole`) VALUES
(1, 'John Doe', 'johndoe@example.com', 'password123', 'user'),
(2, 'reza', 'reza@gmail.com', '123', 'kasir'),
(5, 'reza', 'reza@gmailaaa.com', '12312312312312', ''),
(7, 'reza', 'reza@gmailqaa.com', '12312312312312', ''),
(8, 'reza', 'reza@gmailqaaa.com', '12312312312312', ''),
(10, 'reza', 'reza@gmailqaaaa.com', '12312312312312', ''),
(11, 'reza', 'reza@gmailaqaaaa.com', '12312312312312', 'kasir'),
(12, 'reza', 'reza@gmailaqaaaaa.com', '12312312312312', 'kasir'),
(13, 'rezaa', 'reza@agamailaqaaaaaa.com', '12312312312312', 'kasir'),
(15, 'rezaa', 'reza@agamailaqaaaaaaa.com', '12312312312312', 'kasir');

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `food`
--
ALTER TABLE `food`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `foods`
--
ALTER TABLE `foods`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_foods_deleted_at` (`deleted_at`);

--
-- Indeks untuk tabel `orders`
--
ALTER TABLE `orders`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`);

--
-- Indeks untuk tabel `order_details`
--
ALTER TABLE `order_details`
  ADD PRIMARY KEY (`id`),
  ADD KEY `order_id` (`order_id`),
  ADD KEY `food_id` (`food_id`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD UNIQUE KEY `email_2` (`email`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `food`
--
ALTER TABLE `food`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT untuk tabel `foods`
--
ALTER TABLE `foods`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT untuk tabel `orders`
--
ALTER TABLE `orders`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT untuk tabel `order_details`
--
ALTER TABLE `order_details`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `orders`
--
ALTER TABLE `orders`
  ADD CONSTRAINT `orders_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Ketidakleluasaan untuk tabel `order_details`
--
ALTER TABLE `order_details`
  ADD CONSTRAINT `order_details_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  ADD CONSTRAINT `order_details_ibfk_2` FOREIGN KEY (`food_id`) REFERENCES `food` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
