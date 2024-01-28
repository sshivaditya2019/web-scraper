use std::fmt;

use reqwest::{Error, Response};


pub async fn get_webpage(url: &str) -> Result<Response, Error> {
    let body = reqwest::get(url).await?;
    Ok(body)
}