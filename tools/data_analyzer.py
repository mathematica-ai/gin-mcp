#!/usr/bin/env python3
"""
Data Analyzer Tool for MCP
This tool analyzes data files and provides statistical summaries.
"""

import sys
import json
import pandas as pd
import numpy as np
from pathlib import Path

def analyze_data(file_path):
    """
    Analyze a data file and return statistical summary.
    """
    try:
        # Determine file type and read accordingly
        file_path = Path(file_path)
        
        if file_path.suffix.lower() == '.csv':
            df = pd.read_csv(file_path)
        elif file_path.suffix.lower() in ['.xlsx', '.xls']:
            df = pd.read_excel(file_path)
        elif file_path.suffix.lower() == '.json':
            df = pd.read_json(file_path)
        else:
            return f"Unsupported file format: {file_path.suffix}"
        
        # Generate basic statistics
        stats = {
            "file_info": {
                "filename": file_path.name,
                "file_size": f"{file_path.stat().st_size / 1024:.2f} KB",
                "rows": len(df),
                "columns": len(df.columns)
            },
            "column_info": {},
            "summary_stats": {}
        }
        
        # Analyze each column
        for col in df.columns:
            col_stats = {
                "data_type": str(df[col].dtype),
                "non_null_count": df[col].count(),
                "null_count": df[col].isnull().sum(),
                "null_percentage": f"{(df[col].isnull().sum() / len(df) * 100):.2f}%"
            }
            
            # Add numeric statistics if applicable
            if pd.api.types.is_numeric_dtype(df[col]):
                col_stats.update({
                    "mean": df[col].mean(),
                    "median": df[col].median(),
                    "std": df[col].std(),
                    "min": df[col].min(),
                    "max": df[col].max(),
                    "q25": df[col].quantile(0.25),
                    "q75": df[col].quantile(0.75)
                })
            else:
                # For categorical data, show unique values and most common
                unique_count = df[col].nunique()
                col_stats.update({
                    "unique_values": unique_count,
                    "most_common": df[col].mode().iloc[0] if not df[col].mode().empty else "N/A"
                })
                
                if unique_count <= 10:
                    col_stats["value_counts"] = df[col].value_counts().to_dict()
            
            stats["column_info"][col] = col_stats
        
        # Generate overall summary
        numeric_cols = df.select_dtypes(include=[np.number]).columns
        if len(numeric_cols) > 0:
            stats["summary_stats"]["correlation_matrix"] = df[numeric_cols].corr().to_dict()
        
        return stats
        
    except Exception as e:
        return f"Error analyzing file: {str(e)}"

def main():
    """
    Main function to handle MCP tool execution.
    """
    try:
        # Read input from stdin
        input_data = json.load(sys.stdin)
        
        # Extract arguments
        arguments = input_data.get("arguments", {})
        file_path = arguments.get("file_path")
        
        if not file_path:
            result = {
                "content": [{
                    "type": "text",
                    "text": "Error: file_path argument is required"
                }]
            }
        else:
            # Analyze the data
            analysis_result = analyze_data(file_path)
            
            if isinstance(analysis_result, str) and analysis_result.startswith("Error"):
                result = {
                    "content": [{
                        "type": "text",
                        "text": analysis_result
                    }]
                }
            else:
                result = {
                    "content": [{
                        "type": "text",
                        "text": json.dumps(analysis_result, indent=2, default=str)
                    }]
                }
        
        # Write output to stdout
        json.dump(result, sys.stdout)
        
    except json.JSONDecodeError as e:
        error_result = {
            "content": [{
                "type": "text",
                "text": f"Error parsing input JSON: {str(e)}"
            }]
        }
        json.dump(error_result, sys.stdout)
    except Exception as e:
        error_result = {
            "content": [{
                "type": "text",
                "text": f"Unexpected error: {str(e)}"
            }]
        }
        json.dump(error_result, sys.stdout)

if __name__ == "__main__":
    main() 