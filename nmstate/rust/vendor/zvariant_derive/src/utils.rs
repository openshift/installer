use proc_macro2::TokenStream;
use proc_macro_crate::{crate_name, FoundCrate};
use quote::{format_ident, quote};
use syn::{Attribute, Lit, Meta, Meta::List, MetaList, NestedMeta, Result};

pub fn zvariant_path() -> TokenStream {
    if let Ok(FoundCrate::Name(name)) = crate_name("zvariant") {
        let ident = format_ident!("{}", name);
        quote! { ::#ident }
    } else {
        quote! { ::zvariant }
    }
}

// find the #[@attr_name] attribute in @attrs
fn find_attribute_meta(attrs: &[Attribute], attr_name: &str) -> Result<Option<MetaList>> {
    let meta = match attrs.iter().find(|a| a.path.is_ident(attr_name)) {
        Some(a) => a.parse_meta(),
        _ => return Ok(None),
    };
    match meta? {
        Meta::List(n) => Ok(Some(n)),
        _ => panic!("wrong meta type"),
    }
}

// parse a single meta like: ident = "value"
fn parse_attribute(meta: &NestedMeta) -> (String, String) {
    let meta = match &meta {
        NestedMeta::Meta(m) => m,
        _ => panic!("wrong meta type"),
    };
    let meta = match meta {
        Meta::Path(p) => return (p.get_ident().unwrap().to_string(), "".to_string()),
        Meta::NameValue(n) => n,
        _ => panic!("wrong meta type"),
    };
    let value = match &meta.lit {
        Lit::Str(s) => s.value(),
        _ => panic!("wrong meta type"),
    };

    let ident = match meta.path.get_ident() {
        None => panic!("missing ident"),
        Some(ident) => ident,
    };

    (ident.to_string(), value)
}

#[derive(Debug, PartialEq)]
pub enum ItemAttribute {
    Rename(String),
}

fn parse_item_attribute(meta: &NestedMeta) -> Result<ItemAttribute> {
    let (ident, v) = parse_attribute(meta);

    match ident.as_ref() {
        "rename" => Ok(ItemAttribute::Rename(v)),
        s => panic!("Unknown item meta {}", s),
    }
}

// Parse optional item attributes such as:
// #[zvariant(rename = "MyName")]
pub fn parse_item_attributes(attrs: &[Attribute]) -> Result<Vec<ItemAttribute>> {
    let meta = find_attribute_meta(attrs, "zvariant")?;

    let v = match meta {
        Some(meta) => meta
            .nested
            .iter()
            .map(|m| parse_item_attribute(m).unwrap())
            .collect(),
        None => Vec::new(),
    };

    Ok(v)
}

pub fn get_meta_items(attr: &Attribute) -> Result<Vec<NestedMeta>> {
    if !attr.path.is_ident("zvariant") {
        return Ok(Vec::new());
    }

    match attr.parse_meta() {
        Ok(List(meta)) => Ok(meta.nested.into_iter().collect()),
        _ => panic!("unsupported attribute"),
    }
}
