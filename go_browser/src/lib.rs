#![allow(non_upper_case_globals)]
#![allow(non_camel_case_types)]
#![allow(non_snake_case)]

include!(concat!(env!("OUT_DIR"), "/bindings.rs"));

use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use std::error::Error;
use serde::Serialize;

pub fn make_browser_call(path: &str) -> Result<String, Box<dyn Error>> {
    let c_path = CString::new(path).unwrap();
    println!("Calling Go {}", path);
    let result = unsafe { ExampleBrowser(c_path.as_ptr() as *mut c_char) };
    let c_str = unsafe { CStr::from_ptr(result) };
    let string = c_str.to_str().expect("Error Getting the values from Go");
    match string.is_empty() || string.starts_with("Error") {
        true => Err("Error from Go".into()),
        false => Ok(String::from(string)),
    }
}

#[derive(Serialize)]
pub struct Article {
    pub title: String,
    pub link: String,
    pub time: String,
    pub author: String,
    pub sourcelink: String,
    pub sourcename: String,
}

pub fn get_google_news() -> Vec<Article> {
    let result = unsafe {scrapeGoogleNews() as *mut CArticle};
    let mut articles: Vec<Article> = Vec::new();
    let mut i = 0;
    loop {
        let article = unsafe {&*result.offset(i)};
        println!("Article: {:?}", i);
        if article.Title.is_null() {
            break;
        }
        let title = unsafe {CStr::from_ptr(article.Title).to_str().unwrap()};
        let link = unsafe {CStr::from_ptr(article.Link).to_str().unwrap()};
        let time = unsafe {CStr::from_ptr(article.Time).to_str().unwrap()};
        let author = unsafe {CStr::from_ptr(article.Author).to_str().unwrap()};
        let sourcelink = unsafe {CStr::from_ptr(article.SourceLink).to_str().unwrap()};
        let sourcename = unsafe {CStr::from_ptr(article.SourceName).to_str().unwrap()};
        articles.push(Article {
            title: String::from(title),
            link: String::from(link),
            time: String::from(time),
            author: String::from(author),
            sourcelink: String::from(sourcelink),
            sourcename: String::from(sourcename),
        });
        i += 1;
    }
    unsafe {freeCArticles(result as *mut CArticle, (i-1) as i32)};
    articles
}    
