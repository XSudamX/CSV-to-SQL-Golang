-- phpMyAdmin SQL Dump
-- version 4.9.5deb2
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Dec 28, 2020 at 12:00 PM
-- Server version: 8.0.22-0ubuntu0.20.04.3
-- PHP Version: 7.4.3

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `test`
--

-- --------------------------------------------------------

--
-- Table structure for table `batch_transaction`
--

CREATE TABLE `batch_transaction` (
  `batch_transaction_id` varchar(36) NOT NULL,
  `batch_id` varchar(36) NOT NULL,
  `invoice_id` varchar(36) NOT NULL,
  `payment_method_id` varchar(36) NOT NULL,
  `payment_method_ref_01` varchar(255) NOT NULL,
  `payment_method_ref_02` varchar(255) DEFAULT NULL,
  `amount_gross` double NOT NULL,
  `payment_status` varchar(30) NOT NULL DEFAULT 'PENDING' COMMENT 'PASS/FAIL/PENDING',
  `tried_times` int NOT NULL DEFAULT '0',
  `next_trying_date` datetime DEFAULT NULL,
  `status` varchar(100) NOT NULL DEFAULT 'PENDING' COMMENT 'COMPLETED, PENDING',
  `deleted` tinyint NOT NULL DEFAULT '0',
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `batch_transaction`
--
ALTER TABLE `batch_transaction`
  ADD KEY `fk01_batch_transaction` (`batch_id`);

--
-- Constraints for dumped tables
--

--
-- Constraints for table `batch_transaction`
--
ALTER TABLE `batch_transaction`
  ADD CONSTRAINT `fk01_batch_transaction` FOREIGN KEY (`batch_id`) REFERENCES `batch` (`batch_id`) ON DELETE RESTRICT ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
