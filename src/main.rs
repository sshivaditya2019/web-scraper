use std::error::Error;
use scraper::scraper::get_webpage;
mod scraper;
mod api;


#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    crate::api::api::run_server().await?;
    Ok(())
}

