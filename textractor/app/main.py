from fastapi import FastAPI, HTTPException
from textractor.data.constants import TextractFeatures
from textractor.data.html_linearization_config import HTMLLinearizationConfig

from textractor import Textractor

app = FastAPI()


@app.post("/textract")
async def main(
    file_source: str,
    file_output: str = None,
):
    
    if not file_source.strip():
        raise HTTPException(status_code=400, detail="file_source should not be empty")

    try:
        extractor = Textractor(region_name="ap-south-1")

        document = extractor.start_document_analysis(
            file_source=file_source,
            features=[
                TextractFeatures.FORMS,
                TextractFeatures.LAYOUT,
                TextractFeatures.SIGNATURES,
                TextractFeatures.TABLES,
            ],
            save_image=False,
            s3_output_path=file_output,
        )

        config = HTMLLinearizationConfig(
            table_duplicate_text_in_merged_cells=True
        )
#         config = HTMLLinearizationConfig(
#             figure_layout_prefix="<FIGURE>",
#             figure_layout_suffix="</FIGURE>",
#             header_prefix="<HEADER>",
#             header_suffix="</HEADER>",
#             key_value_layout_prefix="<KEY_VALUE>",
#             key_value_layout_suffix="</KEY_VALUE>",
#             list_element_prefix="<LIST_ELEMENT>",
#             list_element_suffix="</LIST_ELEMENT>",
#             list_layout_prefix="<LIST>",
#             list_layout_suffix="</LIST>",
#             key_prefix="<KEY>",
#             key_suffix="</KEY>",
#             value_prefix="<VALUE>",
#             value_suffix="</VALUE>",
#             page_num_prefix="<PAGE_NUMBER>",
#             page_num_suffix="</PAGE_NUMBER>",
#             section_header_prefix="<SECTION>",
#             section_header_suffix="</SECTION>",
#             selection_element_not_selected="[]",
#             selection_element_selected="[X]",
#             signature_token="[SIGNATURE]",
#             # table_cell_empty_cell_placeholder="",
#             # table_cell_cross_merge_cell_placeholder="",
#             # table_cell_left_merge_cell_placeholder="",
#             # table_cell_merge_cell_placeholder="",
#             # table_cell_header_prefix="",
#             # table_cell_header_suffix="",
#             # table_cell_prefix="",
#             # table_cell_suffix="",
#             # table_cell_top_merge_cell_placeholder="",
#             # table_linearization_format="html",
#             # table_tabulate_format='github'
#         )
        extracted_text = document.get_text(config)
        return {"text": extracted_text}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

