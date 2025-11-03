--
-- PostgreSQL database dump
--

\restrict 5hCPHnx4AaVSnkpVGkUi0tQYPxzk2JvBmFKonTNMjiuwH12JbfKruW2faFM454P

-- Dumped from database version 15.14
-- Dumped by pg_dump version 15.14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: documenttypes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.documenttypes (
    doc_type_id integer NOT NULL,
    doc_type_name character varying(100) NOT NULL
);


--
-- Name: TABLE documenttypes; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.documenttypes IS 'Справочник: Типы сопроводительных документов (УПД, ТОРГ-12)';


--
-- Name: COLUMN documenttypes.doc_type_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.documenttypes.doc_type_id IS 'Код типа документа (PK)';


--
-- Name: COLUMN documenttypes.doc_type_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.documenttypes.doc_type_name IS 'Наименование типа документа';


--
-- Name: documenttypes_doc_type_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.documenttypes_doc_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: documenttypes_doc_type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.documenttypes_doc_type_id_seq OWNED BY public.documenttypes.doc_type_id;


--
-- Name: materialclasses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.materialclasses (
    class_id integer NOT NULL,
    class_name character varying(100) NOT NULL
);


--
-- Name: TABLE materialclasses; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.materialclasses IS 'Справочник: Классы материалов (например, "Строительные материалы")';


--
-- Name: COLUMN materialclasses.class_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materialclasses.class_id IS 'Код класса (PK)';


--
-- Name: COLUMN materialclasses.class_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materialclasses.class_name IS 'Наименование класса';


--
-- Name: materialclasses_class_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.materialclasses_class_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: materialclasses_class_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.materialclasses_class_id_seq OWNED BY public.materialclasses.class_id;


--
-- Name: materialgroups; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.materialgroups (
    group_id integer NOT NULL,
    group_name character varying(100) NOT NULL,
    class_id integer NOT NULL
);


--
-- Name: TABLE materialgroups; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.materialgroups IS 'Справочник: Группы материалов (например, "Сухие смеси")';


--
-- Name: COLUMN materialgroups.group_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materialgroups.group_id IS 'Код группы (PK)';


--
-- Name: COLUMN materialgroups.group_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materialgroups.group_name IS 'Наименование группы';


--
-- Name: COLUMN materialgroups.class_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materialgroups.class_id IS 'Внешний ключ к MaterialClasses';


--
-- Name: materialgroups_group_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.materialgroups_group_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: materialgroups_group_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.materialgroups_group_id_seq OWNED BY public.materialgroups.group_id;


--
-- Name: materials; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.materials (
    material_id bigint NOT NULL,
    material_name character varying(255) NOT NULL,
    group_id integer NOT NULL
);


--
-- Name: TABLE materials; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.materials IS 'Справочник: Конкретные материалы или товары';


--
-- Name: COLUMN materials.material_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materials.material_id IS 'Код материала (PK)';


--
-- Name: COLUMN materials.material_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materials.material_name IS 'Наименование материала';


--
-- Name: COLUMN materials.group_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materials.group_id IS 'Внешний ключ к MaterialGroups';


--
-- Name: materials_material_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.materials_material_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: materials_material_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.materials_material_id_seq OWNED BY public.materials.material_id;


--
-- Name: suppliers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suppliers (
    supplier_id bigint NOT NULL,
    name character varying(255) NOT NULL,
    inn character varying(12) NOT NULL,
    legal_zip_code character varying(10),
    legal_city character varying(100),
    legal_street_address character varying(255),
    bank_zip_code character varying(10),
    bank_city character varying(100),
    bank_street_address character varying(255),
    bank_account character varying(20)
);


--
-- Name: TABLE suppliers; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.suppliers IS 'Справочник: Поставщики материалов';


--
-- Name: COLUMN suppliers.supplier_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.suppliers.supplier_id IS 'Код поставщика (PK)';


--
-- Name: COLUMN suppliers.name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.suppliers.name IS 'Наименование организации';


--
-- Name: COLUMN suppliers.inn; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.suppliers.inn IS 'ИНН (уникальный)';


--
-- Name: suppliers_supplier_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.suppliers_supplier_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: suppliers_supplier_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.suppliers_supplier_id_seq OWNED BY public.suppliers.supplier_id;


--
-- Name: unitsofmeasure; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.unitsofmeasure (
    unit_id integer NOT NULL,
    unit_name character varying(50) NOT NULL,
    abbreviation character varying(10) NOT NULL
);


--
-- Name: TABLE unitsofmeasure; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.unitsofmeasure IS 'Справочник: Единицы измерения (кг, шт, м.п.)';


--
-- Name: COLUMN unitsofmeasure.unit_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.unitsofmeasure.unit_id IS 'Код единицы измерения (PK)';


--
-- Name: COLUMN unitsofmeasure.unit_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.unitsofmeasure.unit_name IS 'Полное наименование';


--
-- Name: COLUMN unitsofmeasure.abbreviation; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.unitsofmeasure.abbreviation IS 'Краткое обозначение';


--
-- Name: unitsofmeasure_unit_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.unitsofmeasure_unit_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: unitsofmeasure_unit_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.unitsofmeasure_unit_id_seq OWNED BY public.unitsofmeasure.unit_id;


--
-- Name: warehousereceipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehousereceipts (
    receipt_id bigint NOT NULL,
    order_number character varying(50) NOT NULL,
    receipt_date date NOT NULL,
    supplier_id bigint NOT NULL,
    balance_account character varying(20) NOT NULL,
    doc_type_id integer NOT NULL,
    document_number character varying(50) NOT NULL,
    material_id bigint NOT NULL,
    material_account character varying(20) NOT NULL,
    unit_id integer NOT NULL,
    quantity numeric(18,3) NOT NULL,
    unit_price numeric(18,2) NOT NULL,
    CONSTRAINT warehousereceipts_quantity_check CHECK ((quantity > (0)::numeric)),
    CONSTRAINT warehousereceipts_unit_price_check CHECK ((unit_price >= (0)::numeric))
);


--
-- Name: TABLE warehousereceipts; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.warehousereceipts IS 'Основная таблица: Приходные ордера (поступления на склад)';


--
-- Name: COLUMN warehousereceipts.receipt_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.warehousereceipts.receipt_id IS 'ID записи о поступлении (PK)';


--
-- Name: COLUMN warehousereceipts.order_number; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.warehousereceipts.order_number IS 'Внутренний номер приходного ордера';


--
-- Name: COLUMN warehousereceipts.receipt_date; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.warehousereceipts.receipt_date IS 'Дата поступления материала на склад';


--
-- Name: COLUMN warehousereceipts.supplier_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.warehousereceipts.supplier_id IS 'Код поставщика (FK)';


--
-- Name: COLUMN warehousereceipts.quantity; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.warehousereceipts.quantity IS 'Количество пришедшего материала (не может быть <= 0)';


--
-- Name: COLUMN warehousereceipts.unit_price; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.warehousereceipts.unit_price IS 'Цена за единицу (не может быть отрицательной)';


--
-- Name: warehousereceipts_receipt_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.warehousereceipts_receipt_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: warehousereceipts_receipt_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.warehousereceipts_receipt_id_seq OWNED BY public.warehousereceipts.receipt_id;


--
-- Name: documenttypes doc_type_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.documenttypes ALTER COLUMN doc_type_id SET DEFAULT nextval('public.documenttypes_doc_type_id_seq'::regclass);


--
-- Name: materialclasses class_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materialclasses ALTER COLUMN class_id SET DEFAULT nextval('public.materialclasses_class_id_seq'::regclass);


--
-- Name: materialgroups group_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materialgroups ALTER COLUMN group_id SET DEFAULT nextval('public.materialgroups_group_id_seq'::regclass);


--
-- Name: materials material_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials ALTER COLUMN material_id SET DEFAULT nextval('public.materials_material_id_seq'::regclass);


--
-- Name: suppliers supplier_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers ALTER COLUMN supplier_id SET DEFAULT nextval('public.suppliers_supplier_id_seq'::regclass);


--
-- Name: unitsofmeasure unit_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unitsofmeasure ALTER COLUMN unit_id SET DEFAULT nextval('public.unitsofmeasure_unit_id_seq'::regclass);


--
-- Name: warehousereceipts receipt_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehousereceipts ALTER COLUMN receipt_id SET DEFAULT nextval('public.warehousereceipts_receipt_id_seq'::regclass);


--
-- Name: documenttypes documenttypes_doc_type_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.documenttypes
    ADD CONSTRAINT documenttypes_doc_type_name_key UNIQUE (doc_type_name);


--
-- Name: documenttypes documenttypes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.documenttypes
    ADD CONSTRAINT documenttypes_pkey PRIMARY KEY (doc_type_id);


--
-- Name: materialclasses materialclasses_class_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materialclasses
    ADD CONSTRAINT materialclasses_class_name_key UNIQUE (class_name);


--
-- Name: materialclasses materialclasses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materialclasses
    ADD CONSTRAINT materialclasses_pkey PRIMARY KEY (class_id);


--
-- Name: materialgroups materialgroups_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materialgroups
    ADD CONSTRAINT materialgroups_pkey PRIMARY KEY (group_id);


--
-- Name: materials materials_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT materials_pkey PRIMARY KEY (material_id);


--
-- Name: suppliers suppliers_inn_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_inn_key UNIQUE (inn);


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (supplier_id);


--
-- Name: unitsofmeasure unitsofmeasure_abbreviation_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unitsofmeasure
    ADD CONSTRAINT unitsofmeasure_abbreviation_key UNIQUE (abbreviation);


--
-- Name: unitsofmeasure unitsofmeasure_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unitsofmeasure
    ADD CONSTRAINT unitsofmeasure_pkey PRIMARY KEY (unit_id);


--
-- Name: unitsofmeasure unitsofmeasure_unit_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unitsofmeasure
    ADD CONSTRAINT unitsofmeasure_unit_name_key UNIQUE (unit_name);


--
-- Name: warehousereceipts warehousereceipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehousereceipts
    ADD CONSTRAINT warehousereceipts_pkey PRIMARY KEY (receipt_id);


--
-- Name: idx_materialgroups_class_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_materialgroups_class_id ON public.materialgroups USING btree (class_id);


--
-- Name: idx_materials_group_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_materials_group_id ON public.materials USING btree (group_id);


--
-- Name: idx_warehousereceipts_material_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_warehousereceipts_material_id ON public.warehousereceipts USING btree (material_id);


--
-- Name: idx_warehousereceipts_receipt_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_warehousereceipts_receipt_date ON public.warehousereceipts USING btree (receipt_date);


--
-- Name: idx_warehousereceipts_supplier_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_warehousereceipts_supplier_id ON public.warehousereceipts USING btree (supplier_id);


--
-- Name: materialgroups fk_materialgroups_class; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materialgroups
    ADD CONSTRAINT fk_materialgroups_class FOREIGN KEY (class_id) REFERENCES public.materialclasses(class_id) ON DELETE RESTRICT;


--
-- Name: materials fk_materials_group; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT fk_materials_group FOREIGN KEY (group_id) REFERENCES public.materialgroups(group_id) ON DELETE RESTRICT;


--
-- Name: warehousereceipts fk_receipts_doctype; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehousereceipts
    ADD CONSTRAINT fk_receipts_doctype FOREIGN KEY (doc_type_id) REFERENCES public.documenttypes(doc_type_id) ON DELETE RESTRICT;


--
-- Name: warehousereceipts fk_receipts_material; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehousereceipts
    ADD CONSTRAINT fk_receipts_material FOREIGN KEY (material_id) REFERENCES public.materials(material_id) ON DELETE RESTRICT;


--
-- Name: warehousereceipts fk_receipts_supplier; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehousereceipts
    ADD CONSTRAINT fk_receipts_supplier FOREIGN KEY (supplier_id) REFERENCES public.suppliers(supplier_id) ON DELETE RESTRICT;


--
-- Name: warehousereceipts fk_receipts_unit; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehousereceipts
    ADD CONSTRAINT fk_receipts_unit FOREIGN KEY (unit_id) REFERENCES public.unitsofmeasure(unit_id) ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

\unrestrict 5hCPHnx4AaVSnkpVGkUi0tQYPxzk2JvBmFKonTNMjiuwH12JbfKruW2faFM454P

