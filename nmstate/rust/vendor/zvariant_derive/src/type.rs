use proc_macro2::TokenStream;
use quote::{quote, ToTokens};
use syn::{self, Attribute, Data, DataEnum, DeriveInput, Fields, Generics, Ident};

use crate::utils::zvariant_path;

pub fn expand_derive(ast: DeriveInput) -> TokenStream {
    let zv = zvariant_path();

    match ast.data {
        Data::Struct(ds) => match ds.fields {
            Fields::Named(_) | Fields::Unnamed(_) => {
                impl_struct(ast.ident, ast.generics, ds.fields, &zv)
            }
            Fields::Unit => impl_unit_struct(ast.ident, ast.generics, &zv),
        },
        Data::Enum(data) => impl_enum(ast.ident, ast.generics, ast.attrs, data, &zv),
        _ => panic!("Only structures and enums supported at the moment"),
    }
}

fn impl_struct(name: Ident, generics: Generics, fields: Fields, zv: &TokenStream) -> TokenStream {
    let (impl_generics, ty_generics, where_clause) = generics.split_for_impl();
    let signature = signature_for_struct(fields, zv);

    quote! {
        impl #impl_generics #zv::Type for #name #ty_generics #where_clause {
            #[inline]
            fn signature() -> #zv::Signature<'static> {
                #signature
            }
        }
    }
}

fn signature_for_struct(fields: Fields, zv: &TokenStream) -> TokenStream {
    let field_types = fields.iter().map(|field| field.ty.to_token_stream());
    let new_type = match fields {
        Fields::Named(_) => false,
        Fields::Unnamed(_) if field_types.len() == 1 => true,
        Fields::Unnamed(_) => false,
        Fields::Unit => panic!("signature_for_struct must not be called for unit fields"),
    };
    if new_type {
        quote! {
            #(
                <#field_types as #zv::Type>::signature()
             )*
        }
    } else {
        quote! {
            let mut s = <::std::string::String as ::std::convert::From<_>>::from("(");
            #(
                s.push_str(<#field_types as #zv::Type>::signature().as_str());
            )*
            s.push_str(")");

            #zv::Signature::from_string_unchecked(s)
        }
    }
}

fn impl_unit_struct(name: Ident, generics: Generics, zv: &TokenStream) -> TokenStream {
    let (impl_generics, ty_generics, where_clause) = generics.split_for_impl();

    quote! {
        impl #impl_generics #zv::Type for #name #ty_generics #where_clause {
            #[inline]
            fn signature() -> #zv::Signature<'static> {
                #zv::Signature::from_static_str_unchecked("")
            }
        }
    }
}

fn impl_enum(
    name: Ident,
    generics: Generics,
    attrs: Vec<Attribute>,
    data: DataEnum,
    zv: &TokenStream,
) -> TokenStream {
    let repr: TokenStream = match attrs.iter().find(|attr| attr.path.is_ident("repr")) {
        Some(repr_attr) => repr_attr
            .parse_args()
            .expect("Failed to parse `#[repr(...)]` attribute"),
        None => quote! { u32 },
    };

    for variant in data.variants {
        // Ensure all variants of the enum are unit type
        match variant.fields {
            Fields::Unit => (),
            _ => panic!("`{}` must be a unit variant", variant.ident),
        }
    }

    let (impl_generics, ty_generics, where_clause) = generics.split_for_impl();

    quote! {
        impl #impl_generics #zv::Type for #name #ty_generics #where_clause {
            #[inline]
            fn signature() -> #zv::Signature<'static> {
                <#repr as #zv::Type>::signature()
            }
        }
    }
}
